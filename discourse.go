package main

import (
	"encoding/json"
	"fmt"
	hooks "github.com/Harvey1717/go-discord-hooks"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/imroc/req"
	"github.com/valyala/fastjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var discourseCategories map[int]string
var newlineRx = regexp.MustCompile("[\n\r ]+")

func getDiscourseCategories(url, token string) (map[int]string, error) {
	res, err := req.Get(url+"categories.json", req.Header{"Api-Key": token})
	if err != nil {
		return nil, err
	}

	var list discourseCategoriesRsp
	err = res.ToJSON(&list)
	if err != nil {
		return nil, err
	}

	categories := make(map[int]string)

	for _, c := range list.CategoryList.Categories {
		if !c.ReadRestricted {
			categories[c.ID] = c.Name
		}
	}

	return categories, nil
}

func getDiscourseTopicSummary(url, token string, topicID int) string {
	res, err := req.Get(fmt.Sprintf("%st/%d.json", url, topicID), req.Header{"Api-Key": token})
	if err != nil {
		return ""
	}

	cooked := fastjson.GetString(res.Bytes(), "post_stream", "posts", "0", "cooked")
	if cooked == "" {
		return ""
	}
	clean := newlineRx.ReplaceAllString(strip.StripTags(cooked), " ")

	return truncateText(clean, 280) + "..."
}

func handleDiscourseWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var u discourseUpdate
	err = json.Unmarshal(body, &u)
	if err != nil {
		return
	}
	if !isValidDiscourseUpdate(u) {
		return
	}

	sendDiscourseNotification(u)
}

func isValidDiscourseUpdate(u discourseUpdate) bool {
	if u.Topic.Archetype != "regular" {
		return false
	}
	if time.Since(u.Topic.CreatedAt) > time.Second*10 {
		return false
	}

	return true
}

func sendDiscourseNotification(u discourseUpdate) {
	category, ok := discourseCategories[u.Topic.CategoryID]
	if !ok {
		return
	}

	e := hooks.NewEmbed()
	e.Title = fmt.Sprintf("[%s] %s", category, u.Topic.Title)
	e.TitleURL = cfg.DiscourseURL + "t/" + u.Topic.Slug
	e.Description = getDiscourseTopicSummary(cfg.DiscourseURL, cfg.DiscourseToken, u.Topic.ID)
	e.Author = hooks.Author{
		Text:    u.Topic.CreatedBy.Name,
		IconURL: cfg.DiscourseURL + strings.ReplaceAll(u.Topic.CreatedBy.AvatarTemplate, "{size}", "128"),
	}
	e.Send(cfg.DiscordWebhook,
		"", "Forum", cfg.DiscourseAvatar)
}

type discourseUpdate struct {
	Topic struct {
		ID         int       `json:"id"`
		Archetype  string    `json:"archetype"`
		Slug       string    `json:"slug"`
		Title      string    `json:"title"`
		CreatedAt  time.Time `json:"created_at"`
		CategoryID int       `json:"category_id"`
		CreatedBy  struct {
			Name           string
			AvatarTemplate string `json:"avatar_template"`
		} `json:"created_by"`
	} `json:"topic"`
}

type discourseCategoriesRsp struct {
	CategoryList struct {
		Categories []discourseCategory `json:"categories"`
	} `json:"category_list"`
}

type discourseCategory struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ReadRestricted bool   `json:"read_restricted"`
}

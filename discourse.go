package main

import (
	"encoding/json"
	"fmt"
	hooks "github.com/Harvey1717/go-discord-hooks"
	"github.com/imroc/req"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var discourseCategories map[int]string

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
	if time.Since(u.Topic.CreatedAt) > time.Second*30 {
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
	e.Author = hooks.Author{
		Text:    u.Topic.CreatedBy.Name,
		IconURL: cfg.DiscourseURL + strings.ReplaceAll(u.Topic.CreatedBy.AvatarTemplate, "{size}", "128"),
	}
	e.Send(cfg.DiscordWebhook,
		"", "Forum", cfg.DiscourseAvatar)
}

type discourseUpdate struct {
	Topic struct {
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

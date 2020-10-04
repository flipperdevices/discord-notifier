package main

import (
	"encoding/json"
	"fmt"
	hooks "github.com/Harvey1717/go-discord-hooks"
	"io/ioutil"
	"net/http"
)

func handleGhostWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var u ghostUpdate
	err = json.Unmarshal(body, &u)
	if err != nil {
		return
	}

	sendGhostNotification(u)
}

func sendGhostNotification(u ghostUpdate) {
	p := u.Post.Current
	if p.Tag.Name != "" {
		p.Title = fmt.Sprintf("[%s] %s", p.Tag.Name, p.Title)
	}

	e := hooks.NewEmbed()

	e.Title = p.Title
	e.TitleURL = p.URL
	e.Description = p.Excerpt + "..."
	e.Author = hooks.Author{
		Text:    p.Author.Name,
		IconURL: p.Author.Image,
	}
	e.Send(cfg.DiscordWebhook,
		"", "Blog", cfg.GhostAvatar)
}

type ghostUpdate struct {
	Post struct {
		Current struct {
			Title  string `json:"title"`
			URL    string `json:"url"`
			Author struct {
				Name  string `json:"name"`
				Image string `json:"profile_image"`
			} `json:"primary_author"`
			Tag struct {
				Name string `json:"name"`
			} `json:"primary_tag"`
			Excerpt string `json:"excerpt"`
		} `json:"current"`
	} `json:"post"`
}

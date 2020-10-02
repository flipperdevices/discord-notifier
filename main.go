package main

import (
	"github.com/caarlos0/env"
	"log"
	"net/http"
)

var cfg config

func main() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Config", err)
	}

	var err error

	discourseCategories, err = getDiscourseCategories(cfg.DiscourseURL, cfg.DiscourseToken)
	if err != nil {
		log.Fatalln("Discourse", err)
	}

	log.Println("Server started")
	http.HandleFunc("/discourse", handleDiscourseWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

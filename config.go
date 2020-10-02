package main

type config struct {
	DiscourseURL    string `env:"DISCOURSE_URL,required"`
	DiscourseToken  string `env:"DISCOURSE_TOKEN,required"`
	DiscourseAvatar string `env:"DISCOURSE_AVATAR"`

	DiscordWebhook string `env:"DISCORD_WEBHOOK,required"`
}

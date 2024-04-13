package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	WebhookUrlsFormat string
	WebhookUrlSecret  string
	ElasticSearchHost string
}

var Configs Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	Configs = Config{
		ServerPort:        os.Getenv("SERVER_PORT"),
		WebhookUrlsFormat: os.Getenv("WEBHOOK_URLS_FORMAT"),
		WebhookUrlSecret:  os.Getenv("WEBHOOK_URL_SECRET"),
		ElasticSearchHost: os.Getenv("ELASTICSEARCH_HOST"),
	}
}

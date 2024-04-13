package main

import (
	"BotMother/config"
	"BotMother/db"
	"BotMother/logger"
	"BotMother/web"
)

func main() {
	// Load configurations
	config.LoadConfig()

	// Initialize logger
	logger.Init()

	// Initialize database
	db.Init()

	// Available bots
	availableBots := []string{"first", "second"}

	// Start the web server
	web.StartServer(availableBots)
}

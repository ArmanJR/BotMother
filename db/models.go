package db

import (
	"gorm.io/gorm"
)

// User represents a user of the Telegram bot in the database.
type User struct {
	gorm.Model
	BotID   string `gorm:"index"` // Index this column for faster queries by BotID
	StateID string
}

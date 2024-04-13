package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// Init initializes the database connection.
func Init() {
	var err error
	dsn := "host=localhost user=youruser dbname=yourdb sslmode=disable password=yourpassword"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&User{})
}

// CreateUser creates a new user in the database.
func CreateUser(user *User) error {
	result := DB.Create(user)
	return result.Error
}

// UpdateUserState updates the state of an existing user.
func UpdateUserState(userID uint, newState string) error {
	result := DB.Model(&User{}).Where("id = ?", userID).Update("state_id", newState)
	return result.Error
}

// GetUserByBotID retrieves a user by their BotID.
func GetUserByBotID(botID string) (*User, error) {
	var user User
	result := DB.Where("bot_id = ?", botID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

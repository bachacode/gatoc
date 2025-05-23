package database

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type WelcomeRole struct {
	gorm.Model
	GuildID string
	RoleID  string
	UserID  *string
}

func New(cfg *config.DbConfig) (*gorm.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&WelcomeRole{})

	if err != nil {
		return err
	}

	return nil
}

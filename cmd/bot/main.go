package main

import (
	"log"
	"os"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bachacode/go-discord-bot/internal/config"
	"github.com/bachacode/go-discord-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags)
	cfg := config.LoadConfig()

	b, err := bot.New(cfg)
	if err != nil {
		logger.Fatalf("Failed to create bot: %v", err)
		return
	}
	_, err = database.New(cfg.DbConfig)
	if err != nil {
		logger.Fatalf("Failed to connect to db: %v", err)
		return
	}

	err = b.Setup(discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates)
	if err != nil {
		logger.Fatalf("Failed to setup bot: %v", err)
		return
	}

	err = b.Run()
	if err != nil {
		logger.Fatalln(err)
		return
	}
}

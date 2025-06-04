package main

import (
	"log"
	"os"

	"github.com/bachacode/go-discord-bot/internal/bot"
	_ "github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bachacode/go-discord-bot/internal/config"
	"github.com/bachacode/go-discord-bot/internal/database"
	_ "github.com/bachacode/go-discord-bot/internal/events"
	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags|log.Lshortfile)
	cfg := config.LoadConfig()

	db, err := database.New(cfg.DbConfig)
	if err != nil {
		logger.Fatalf("ERROR: Failed to connect to db: %v\n", err)
		return
	}

	if err := database.Migrate(db); err != nil {
		logger.Fatalf("ERROR: Failed to migrate tables to db: %v\n", err)
	}

	bb := bot.NewBotBuilder(cfg.BotConfig)
	bb.WithDatabase(db)
	bb.WithIntents(
		discordgo.IntentsGuilds |
			discordgo.IntentsGuildVoiceStates |
			discordgo.IntentsMessageContent |
			discordgo.IntentGuildMessages,
	)
	bb.WithLogger(logger)
	bot, err := bb.Build()
	if err != nil {
		logger.Fatalf("ERROR: Failed to create bot instance: %v\n", err)
		return
	}

	bot.Setup()

	if err := bot.Run(); err != nil {
		logger.Fatalf("ERROR: Failed to run bot instance: %v\n", err)
		return
	}
}

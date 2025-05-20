package main

import (
	"log"
	"os"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	logger := log.New(os.Stdout, "[MAIN] ", log.LstdFlags)
	env := os.Getenv("ENV")

	var token, appID, guildID string
	if env == "development" {
		token = os.Getenv("TEST_TOKEN")
		appID = os.Getenv("TEST_CLIENT_ID")
	} else {
		token = os.Getenv("TOKEN")
		appID = os.Getenv("CLIENT_ID")
	}
	guildID = os.Getenv("TEST_GUILD_ID")

	b, err := bot.New(token, appID, guildID)

	if err != nil {
		logger.Fatalf("Failed to create bot: %v", err)
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

package main

import (
	"fmt"
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
	token := os.Getenv("TEST_TOKEN")
	appID := os.Getenv("TEST_CLIENT_ID")
	guildID := os.Getenv("TEST_GUILD_ID")

	b, err := bot.New(token, appID, guildID)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = b.Setup(discordgo.IntentsGuilds | discordgo.IntentsGuildVoiceStates)
	if err != nil {
		fmt.Println("error setting up the bot,", err)
		return
	}

	err = b.Run()
	if err != nil {
		fmt.Println("error running the bot,", err)
		return
	}
}

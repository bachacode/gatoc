package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bachacode/go-discord-bot/internal/bot"
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
	appId := os.Getenv("TEST_CLIENT_ID")
	guildId := os.Getenv("TEST_GUILD_ID")

	err := bot.Init(token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = bot.SetEventHandlers()
	if err != nil {
		fmt.Println("error setting event handlers,", err)
		return
	}

	err = bot.RegisterCommands(appId, guildId)
	if err != nil {
		fmt.Println("error registering commands,", err)
		return
	}

	err = bot.Start()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = bot.UnregisterCommands(guildId)
	if err != nil {
		fmt.Println("error unregistering commands,", err)
		return
	}

	err = bot.Close()
	if err != nil {
		fmt.Println("error closing connection,", err)
		return
	}
}

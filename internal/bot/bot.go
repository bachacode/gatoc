package bot

import (
	"errors"
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bachacode/go-discord-bot/internal/events"
	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session
var errSessionNotInitialized = errors.New("bot: session not initialized; call Init() first")

func Init(token string) error {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages
	session = s
	return nil
}

func SetEventHandlers() error {
	if session == nil {
		return errSessionNotInitialized
	}

	session.AddHandler(events.ReadyHandler)
	session.AddHandler(events.MessageCreateHandler)
	session.AddHandler(events.InteractionCreateHandler)
	return nil
}

func RegisterCommands(appId string, guildId string) error {
	for _, cmd := range commands.Commands {
		_, err := session.ApplicationCommandCreate(appId, guildId, cmd.Metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnregisterCommands(guildId string) error {
	appId := session.State.User.ID

	commands, err := session.ApplicationCommands(appId, guildId)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		err := session.ApplicationCommandDelete(appId, guildId, cmd.ID)
		if err != nil {
			fmt.Printf("Error deleting command %s: %v\n", cmd.Name, err)
		}
	}

	return nil
}

func Start() error {
	if session == nil {
		return errSessionNotInitialized
	}
	return session.Open()
}

func Close() error {
	if session == nil {
		return errSessionNotInitialized
	}

	return session.Close()
}

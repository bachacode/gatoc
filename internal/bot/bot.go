package bot

import (
	"errors"

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

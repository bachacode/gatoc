package bot

import (
	"errors"

	"github.com/bachacode/go-discord-bot/internal/events"
	"github.com/bwmarrin/discordgo"
)

var session *discordgo.Session
var errSessionNotInitialized = errors.New("bot: session not initialized; call Init() first")

// Init initializes the Discord session and sets up basic state.
// This should always be called first.
func Init(token string) error {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	s.Identify.Intents = discordgo.IntentsGuildMessages
	session = s
	return nil
}

// SetEventHandlers sets all the handlers needed by the bot.
// It should be called after Init, otherwise returns an error.
func SetEventHandlers() error {
	if session == nil {
		return errSessionNotInitialized
	}

	session.AddHandler(events.ReadyHandler)
	session.AddHandler(events.MessageCreateHandler)
	return nil
}

// Start opens the Discord connection.
func Start() error {
	if session == nil {
		return errSessionNotInitialized
	}

	return session.Open()
}

// Close shuts down the Discord connection.
func Close() error {
	if session == nil {
		return errSessionNotInitialized
	}

	return session.Close()
}

package bot

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bachacode/go-discord-bot/internal/events"
	"github.com/bwmarrin/discordgo"
)

type bot struct {
	session  *discordgo.Session
	appID    string
	guildID  string
	intents  discordgo.Intent
	commands map[string]commands.SlashCommand
	events   []events.EventHandler
}

var errSessionNotInitialized = errors.New("bot: session not initialized; call Init() first")

func New(token string, appID string, guildID string) (*bot, error) {

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	bot := &bot{
		session:  s,
		appID:    appID,
		guildID:  guildID,
		commands: commands.All(),
		events:   events.All(),
	}
	return bot, nil
}

func (b *bot) Setup(intents discordgo.Intent) error {
	if b.session == nil {
		return errSessionNotInitialized
	}

	b.session.Identify.Intents = intents
	b.setEventHandlers()
	b.registerCommands()
	return nil
}

func (b *bot) setEventHandlers() {
	for _, event := range b.events {
		if event.Once {
			b.session.AddHandlerOnce(event.Handler)
		} else {
			b.session.AddHandler(event.Handler)
		}
	}
}

func (b *bot) registerCommands() error {
	for _, cmd := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.appID, b.guildID, cmd.Metadata)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *bot) UnregisterCommands() error {
	commands, err := b.session.ApplicationCommands(b.appID, b.guildID)
	if err != nil {
		return err
	}

	for _, cmd := range commands {
		err := b.session.ApplicationCommandDelete(b.appID, b.guildID, cmd.ID)
		if err != nil {
			fmt.Printf("Error deleting command %s: %v\n", cmd.Name, err)
		}
	}

	return nil
}

func (b *bot) Run() error {
	if b.session == nil {
		return errSessionNotInitialized
	}

	b.session.Open()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err := b.UnregisterCommands()
	if err != nil {
		return err
	}
	b.session.Close()

	return nil
}

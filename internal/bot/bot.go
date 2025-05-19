package bot

import (
	"fmt"
	"log"
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
	events   []events.Event
	logger   *log.Logger
}

func New(token string, appID string, guildID string) (*bot, error) {

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "[BOT] ", log.LstdFlags|log.Lshortfile)

	bot := &bot{
		session:  s,
		appID:    appID,
		guildID:  guildID,
		commands: commands.All(),
		events:   events.All(),
		logger:   logger,
	}
	return bot, nil
}

func (b *bot) Setup(intents discordgo.Intent) error {
	b.session.Identify.Intents = intents

	b.logger.Println("Setting bot events...")
	b.setEventHandlers()

	b.logger.Println("Registering commands...")
	err := b.registerCommands()
	if err != nil {
		return fmt.Errorf("Error registering commands: %v", err)
	}

	return nil
}

func (b *bot) setEventHandlers() {
	for _, event := range b.events {
		if event.Once {
			b.session.AddHandlerOnce(event.Handler)
			b.logger.Printf("Registered event: %s as once event", event.Name)
		} else {
			b.session.AddHandler(event.Handler)
			b.logger.Printf("Registered event: %s as normal event", event.Name)
		}
	}
	b.logger.Printf("All events were register successfully!")
}

func (b *bot) registerCommands() error {
	for _, cmd := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.appID, b.guildID, cmd.Metadata)
		if err != nil {
			return fmt.Errorf("Failed to register command: %s\n%v", cmd.Metadata.Name, err)
		}
		b.logger.Printf("Registered command: %s\n", cmd.Metadata.Name)
	}
	b.logger.Printf("All commands were register successfully!")
	return nil
}

func (b *bot) UnregisterCommands() error {
	for _, cmd := range b.commands {
		err := b.session.ApplicationCommandDelete(b.appID, b.guildID, cmd.Metadata.ID)
		if err != nil {
			return fmt.Errorf("Failed to delete command: %s\n%v", cmd.Metadata.Name, err)
		}
	}
	b.logger.Printf("All commands were unregister successfully!")
	return nil
}

func (b *bot) Run() error {

	b.logger.Println("Starting bot session...")
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error starting bot session: %v", err)
	}

	b.logger.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.logger.Println("Unregistering commands...")
	if err := b.UnregisterCommands(); err != nil {
		return fmt.Errorf("Error unregistering commands: %v", err)
	}

	b.logger.Println("Closing bot session...")
	b.session.Close()
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error closing bot session: %v", err)
	}

	return nil
}

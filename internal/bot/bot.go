package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bachacode/go-discord-bot/internal/config"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type BotContext struct {
	*config.BotConfig
	DB     *gorm.DB
	Logger *log.Logger
}

type bot struct {
	session  *discordgo.Session
	intents  discordgo.Intent
	commands map[string]SlashCommand
	events   []Event
	*BotContext
}

func (b *bot) Setup() {
	b.Logger.Println("INFO: Setting bot events...")
	b.setEventHandlers()

	b.Logger.Println("INFO: Registering commands...")
	b.registerCommands()
}

func (b *bot) setEventHandlers() {
	for _, event := range b.events {
		if event.Once {
			b.session.AddHandlerOnce(event.Handler(b.BotContext))
			b.Logger.Printf("INFO: Registered event: %s as once event\n", event.Name)
		} else {
			b.session.AddHandler(event.Handler(b.BotContext))
			b.Logger.Printf("INFO: Registered event: %s as normal event\n", event.Name)
		}
	}
	b.Logger.Println("INFO: All events were register successfully!")
}

func (b *bot) registerCommands() {
	total := len(b.commands)
	count := 0
	for _, cmd := range b.commands {
		_, err := b.session.ApplicationCommandCreate(b.ClientID, b.GuildID, cmd.Metadata)

		if err != nil {
			b.Logger.Printf("WARN: Failed to register command %s: %v\n", cmd.Metadata.Name, err)
		} else {
			count++
			b.Logger.Printf("INFO: Registered command: %s\n", cmd.Metadata.Name)
		}
	}
	b.Logger.Printf("INFO: %d commands of %d were register successfully!", count, total)
}

func (b *bot) UnregisterCommands() {
	commands, err := b.session.ApplicationCommands(b.ClientID, b.GuildID)

	if err != nil {
		b.Logger.Printf("WARN: Failed to fetch applications commands: %v", err)
		b.Logger.Println("WARN: Skipping commands removal...")
		return
	}

	total := len(commands)
	count := 0

	for _, cmd := range commands {
		err := b.session.ApplicationCommandDelete(b.ClientID, b.GuildID, cmd.ID)
		if err != nil {
			b.Logger.Printf("WARN: Failed to unregister command %s: %v\n", cmd.Name, err)
		} else {
			count++
			b.Logger.Printf("INFO: Unregistered command: %s\n", cmd.Name)
		}
	}
	b.Logger.Printf("%d commands of %d were unregister successfully!\n", count, total)
}

func (b *bot) Run() error {

	b.session.Identify.Intents = b.intents
	b.Logger.Println("INFO: Starting bot session...")
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error starting bot session: %v", err)
	}

	b.Logger.Println("INFO: Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	b.Logger.Println("INFO: Unregistering commands...")
	b.UnregisterCommands()

	b.Logger.Println("INFO: Closing bot session...")
	b.session.Close()
	if err := b.session.Open(); err != nil {
		return fmt.Errorf("Error closing bot session: %v", err)
	}

	return nil
}

package bot

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bachacode/go-discord-bot/internal/config"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

type BotBuilder struct {
	// Mandatory configuration
	cfg *config.BotConfig

	// Optional dependencies/configurations
	db             *gorm.DB
	logger         *log.Logger
	intents        discordgo.Intent
	customCommands map[string]commands.SlashCommand
	customEvents   []Event
}

func NewBotBuilder(cfg *config.BotConfig) *BotBuilder {
	return &BotBuilder{
		cfg:            cfg,
		logger:         log.New(os.Stdout, "[DEFAULT_BOT] ", log.LstdFlags|log.Lshortfile),
		intents:        discordgo.IntentGuildMessages,
		customCommands: make(map[string]commands.SlashCommand),
		customEvents:   make([]Event, 0),
	}
}

func (bb *BotBuilder) WithDatabase(db *gorm.DB) *BotBuilder {
	bb.db = db
	return bb
}

func (bb *BotBuilder) WithLogger(logger *log.Logger) *BotBuilder {
	bb.logger = logger
	return bb
}

func (bb *BotBuilder) WithIntents(intents discordgo.Intent) *BotBuilder {
	bb.intents = intents
	return bb
}

func (bb *BotBuilder) WithCustomCommands(cmds map[string]commands.SlashCommand) *BotBuilder {
	bb.customCommands = cmds
	return bb
}

func (bb *BotBuilder) WithCustomEvents(evts []Event) *BotBuilder {
	bb.customEvents = evts
	return bb
}

func (bb *BotBuilder) Build() (*bot, error) {
	// Validate that base config and token exist
	if bb.cfg == nil {
		return nil, errors.New("bot configuration is required")
	}
	if bb.cfg.Token == "" {
		return nil, errors.New("Discord bot token is required in the configuration")
	}

	// Create new discord session
	s, err := discordgo.New("Bot " + bb.cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	// Create bot context
	botCtx := &BotContext{
		BotConfig: bb.cfg,
		DB:        bb.db,
		Logger:    bb.logger,
	}

	// Create new bot struct
	b := &bot{
		session:    s,
		intents:    bb.intents,
		BotContext: botCtx,
	}

	// Create commands
	if len(bb.customCommands) > 0 {
		b.commands = bb.customCommands
		botCtx.Logger.Println("INFO: Using custom commands")
	} else {
		b.commands = commands.All()
		botCtx.Logger.Println("INFO: Using default commands")
	}

	// Create events
	if len(bb.customEvents) > 0 {
		b.events = bb.customEvents
		botCtx.Logger.Println("INFO: Using custom events")
	} else {
		b.events = eventRegistry
		botCtx.Logger.Println("INFO: Using default events")
	}

	botCtx.Logger.Println("INFO: Bot instance successfully built")
	return b, nil
}

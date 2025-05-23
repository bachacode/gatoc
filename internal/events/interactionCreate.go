package events

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(interactionCreate)
}

var interactionCreate bot.Event = bot.Event{
	Name: "Interaction Create / Slash Command Handling",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type != discordgo.InteractionApplicationCommand {
				return
			}

			cmd := bot.GetCommand(i.ApplicationCommandData().Name)
			if err := cmd.Handler(s, i, ctx); err != nil {
				fmt.Printf("Failed to run interaction: %v\n", err)
			}
			return
		}
	},
}

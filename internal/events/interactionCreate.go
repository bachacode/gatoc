package events

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

func init() {
	register(interactionCreate)
}

var interactionCreate Event = Event{
	Name: "Interaction Create / Slash Command Handling",
	Once: false,
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		cmd := commands.Get(i.ApplicationCommandData().Name)
		if err := cmd.Handler(s, i); err != nil {
			fmt.Printf("Failed to run interaction: %v\n", err)
		}
		return
	},
}

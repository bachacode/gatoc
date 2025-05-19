package events

import (
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
		cmd.Handler(s, i)
		return
	},
}

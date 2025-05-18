package events

import (
	"github.com/bachacode/go-discord-bot/internal/commands"
	"github.com/bwmarrin/discordgo"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	for _, cmd := range commands.Commands {
		if i.ApplicationCommandData().Name == cmd.Metadata.Name {
			cmd.Handler(s, i)
			return
		}
	}
}

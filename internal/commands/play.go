package commands

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	register(play.Metadata.Name, play)
}

var play SlashCommand = SlashCommand{
	Metadata: &discordgo.ApplicationCommand{}}

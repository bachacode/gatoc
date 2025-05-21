package events

import (
	"log"

	"github.com/bachacode/go-discord-bot/internal/config"
	"github.com/bwmarrin/discordgo"
)

func init() {
	register(ready)
}

var ready Event = Event{
	Name: "Ready",
	Once: true,
	Handler: func(cfg *config.BotConfig) interface{} {
		return func(s *discordgo.Session, r *discordgo.Ready) {
			log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		}
	},
}

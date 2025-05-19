package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(ready)
}

var ready EventHandler = EventHandler{
	Name: "Ready",
	Once: true,
	Handler: func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	},
}

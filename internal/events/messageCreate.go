package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(messageCreate)
}

var messageCreate EventHandler = EventHandler{
	Name: "Message Create",
	Once: false,
	Handler: func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID || m.Author.Bot {
			return
		}

		if m.Content != "ping" {
			return
		}

		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("error creating channel:", err)
			return
		}

		_, err = s.ChannelMessageSendReply(channel.ID, "Pong!", &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		})
		if err != nil {
			fmt.Println("error sending a reply message:", err)
			s.ChannelMessageSend(
				m.ChannelID,
				"Ha ocurrido un error al responder a tu mensaje",
			)
		}
	},
}

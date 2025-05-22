package events

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(messageCreate)
}

var messageCreate bot.Event = bot.Event{
	Name: "Message Create",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		}
	},
}

package events

import (
	"strings"

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

			channelID := m.ChannelID
			messages, err := s.ChannelMessages(channelID, 2, "", "", "")

			if err != nil {
				ctx.Logger.Printf("Failed to get message history from channel: %v", err)
				return
			}

			count := 0

			if strings.ToLower(messages[0].Content) == strings.ToLower(messages[1].Content) &&
				(messages[0].Author.GlobalName != messages[1].Author.GlobalName) {
				count++
			} else {
				count = 0
			}

			if count >= 2 {
				count = 0
				_, err := s.ChannelMessageSend(channelID, messages[len(messages)-1].Content)
				if err != nil {
					ctx.Logger.Printf("Failed to send message: %v", err)
					return
				}
			}
		}
	},
}

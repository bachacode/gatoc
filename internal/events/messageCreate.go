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
			channel, err := s.Channel(channelID)

			if err != nil {
				ctx.Logger.Printf("Failed to get channel from the message: %v", err)
				return
			}

			if len(channel.Messages) >= 3 {
				allSame := true
				lastMessages := channel.Messages[len(channel.Messages)-3:]

				for _, message := range lastMessages {
					if strings.ToLower(message.Content) != strings.ToLower(lastMessages[0].Content) {
						allSame = false
					}
				}

				if allSame {
					_, err := s.ChannelMessageSend(channelID, lastMessages[len(lastMessages)-1].Content)

					if err != nil {
						ctx.Logger.Printf("Failed to send message: %v", err)
						return
					}
				}
			}
		}
	},
}

package events

import (
	"strings"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(messageReactionAdd)
}

var messageReactionAdd bot.Event = bot.Event{
	Name: "Message Reaction Add",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
			user, err := s.User(m.UserID)
			if err != nil {
				ctx.Logger.Printf("Error fetching user: %v", err)
				return
			}

			if user.Bot {
				return
			}

			targetEmojiID := "957421664738639872"
			channelID := m.ChannelID
			messageID := m.MessageID
			botID := ctx.ClientID
			if m.Emoji.ID != targetEmojiID {
				return
			}

			msg, err := s.ChannelMessage(channelID, messageID)
			if err != nil {
				ctx.Logger.Printf("Failed to get message: %v", err)
				return
			}

			if strings.Contains(msg.Content, "[Fix]") && msg.Author.ID == botID {
				newContent := msg.Content + " "
				_, err = s.ChannelMessageEdit(channelID, messageID, newContent)
				if err != nil {
					ctx.Logger.Printf("Error editing message: %v", err)
					return
				}
			}
		}
	},
}

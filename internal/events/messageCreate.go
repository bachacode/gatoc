package events

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(messageCreate)
}

var messageCount int = 0

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

			// Repeat message if last 3 messages are the same
			if areMessagesRepeated(messages, 2) {
				messageCount = 0
				_, err := s.ChannelMessageSend(channelID, messages[len(messages)-1].Content)
				if err != nil {
					ctx.Logger.Printf("Failed to send message: %v", err)
					return
				}
			}

			var messageEdit *discordgo.MessageEdit
			var fixedMessage *string
			isFxtwitter := strings.Contains(m.Content, "fxtwitter.com")
			isTwitterOrX := strings.Contains(m.Content, "twitter.com") || strings.Contains(m.Content, "x.com")
			hasStatusPath := strings.Contains(m.Content, "/status/")

			if !isFxtwitter && isTwitterOrX && hasStatusPath {
				fixedMessage = fixTwitterEmbed(m)
				messageEdit = &discordgo.MessageEdit{
					ID:      m.ID,
					Channel: m.ChannelID,
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				}
			}

			isVxreddit := strings.Contains(m.Content, "vxreddit.com")
			isReddit := strings.Contains(m.Content, "reddit.com")
			hasCommentsPath := strings.Contains(m.Content, "/comments/")
			if !isVxreddit && isReddit && hasCommentsPath {
				fixedMessage = fixRedditEmbed(m)
				messageEdit = &discordgo.MessageEdit{
					ID:      m.ID,
					Channel: m.ChannelID,
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				}
			}

			isDDinstagram := strings.Contains(m.Content, "ddinstagram.com")
			isInstagram := strings.Contains(m.Content, "instagram.com")
			hasPPath := strings.Contains(m.Content, "/p/")
			if !isDDinstagram && isInstagram && hasPPath {
				fixedMessage = fixInstagramEmbed(m)
				messageEdit = &discordgo.MessageEdit{
					ID:      m.ID,
					Channel: m.ChannelID,
					Flags:   discordgo.MessageFlagsSuppressEmbeds,
				}
			}

			if messageEdit != nil {
				// Supress embeds
				if _, err := s.ChannelMessageEditComplex(messageEdit); err != nil {
					ctx.Logger.Printf("Failed to supress embeds from previous message: %v", err)
					return
				}

				if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
					Content: *fixedMessage,
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{},
					},
				}); err != nil {
					ctx.Logger.Printf("Failed to send embedded message: %v", err)
					return
				}
			}

		}
	},
}

func areMessagesRepeated(messages []*discordgo.Message, max int) bool {
	if strings.ToLower(messages[0].Content) == strings.ToLower(messages[1].Content) &&
		(messages[0].Author.GlobalName != messages[1].Author.GlobalName) {
		messageCount++
	} else {
		messageCount = 0
	}

	return messageCount >= max
}

func fixTwitterEmbed(m *discordgo.MessageCreate) *string {
	// Regex pattern to match a URL
	re := regexp.MustCompile(`https?://[^\s]+`)
	// Find the first URL match
	url := re.FindString(m.Content)

	tweet := "<" + url + ">"
	authorName := strings.Split(strings.Split(url, "/status")[0], ".com/")[1]
	author := "<" + strings.Split(url, "/status")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	fxtwitterURL := strings.Replace(url, "twitter.com", "fxtwitter.com", 1)
	fxtwitterURL = strings.Replace(fxtwitterURL, "x.com", "fxtwitter.com", 1)

	s := fmt.Sprintf("[Tweet](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", tweet, authorName, author, fxtwitterURL, mention)
	return &s
}

func fixRedditEmbed(m *discordgo.MessageCreate) *string {
	// Regex pattern to match a URL
	re := regexp.MustCompile(`https?://[^\s]+`)
	// Find the first URL match
	url := re.FindString(m.Content)

	redditURL := "<" + url + ">"
	authorName := strings.Split(strings.Split(url, "/comments")[0], "r/")[1]
	author := "<" + strings.Split(url, "/comments")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	vxredditURL := strings.Replace(url, "reddit.com", "vxreddit.com", 1)

	s := fmt.Sprintf("[Reddit](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", redditURL, authorName, author, vxredditURL, mention)
	return &s
}

func fixInstagramEmbed(m *discordgo.MessageCreate) *string {
	// Regex pattern to match a URL
	re := regexp.MustCompile(`https?://[^\s]+`)
	// Find the first URL match
	url := re.FindString(m.Content)

	instagramURL := "<" + url + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	ddinstagramURL := strings.Replace(url, "instagram.com", "ddinstagram.com", 1)

	s := fmt.Sprintf("[Instagram](%s) • [Fix](%s) • Enviado por %s ", instagramURL, ddinstagramURL, mention)
	return &s
}

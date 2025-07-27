package events

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(messageCreate)
}

var messageCount = 0

var messageCreate bot.Event = bot.Event{
	Name: "Message Create",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, m *discordgo.MessageCreate) {
			if m.Author.ID == s.State.User.ID || m.Author.Bot {
				return
			}

			channelID := m.ChannelID
			maxMessages := 3

			if err := handleRepeated(channelID, maxMessages, s); err != nil {
				ctx.Logger.Printf("Failed to repeat messages: %v", err)
			}

			if err := handleURLEmbed(s, m); err != nil {
				ctx.Logger.Printf("Failed to fix message embed: %v", err)
			}

		}
	},
}

func handleRepeated(channelID string, max int, s *discordgo.Session) error {
	messages, err := s.ChannelMessages(channelID, 2, "", "", "")
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return fmt.Errorf("Failed to get enough messages from message history")
	}

	isSameMessage := strings.ToLower(messages[0].Content) == strings.ToLower(messages[1].Content)
	isDifferentAuthor := messages[0].Author.GlobalName != messages[1].Author.GlobalName
	if isSameMessage && isDifferentAuthor {
		messageCount++
	} else {
		messageCount = 1
	}

	if messageCount >= max {
		messageCount = 0
		_, err := s.ChannelMessageSend(channelID, messages[len(messages)-1].Content)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleURLEmbed(s *discordgo.Session, m *discordgo.MessageCreate) error {
	url, err := url.ParseRequestURI(m.Content)
	if err != nil {
		return nil
	}

	trimmedHost := strings.ToLower(url.Host)
	if strings.HasPrefix(trimmedHost, "www.") {
		trimmedHost = strings.TrimPrefix(trimmedHost, "www.")
	}

	fixableHosts := map[string]func(m *discordgo.MessageCreate) string{
		"twitter.com":   fixTwitterEmbed,
		"x.com":         fixTwitterEmbed,
		"reddit.com":    fixRedditEmbed,
		"instagram.com": fixInstagramEmbed,
	}

	if handler, ok := fixableHosts[trimmedHost]; ok {
		if (trimmedHost == "twitter.com" || trimmedHost == "x.com") && !strings.Contains(url.Path, "/status") {
			return nil
		}

		fixedEmbedMessageContent := handler(m)

		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      m.ID,
				Channel: m.ChannelID,
				Flags:   discordgo.MessageFlagsSuppressEmbeds,
			})
		}

		if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content: fixedEmbedMessageContent,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Parse: []discordgo.AllowedMentionType{},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func fixUrlEmbed(m *discordgo.MessageCreate) (string, string) {
	re := regexp.MustCompile(`https?://[^\s]+`)

	url := re.FindString(m.Content)
	originalSupressedUrl := "<" + url + ">"

	return url, originalSupressedUrl
}

func fixTwitterEmbed(m *discordgo.MessageCreate) string {
	url, fixedUrl := fixUrlEmbed(m)

	authorName := strings.Split(strings.Split(url, "/status")[0], ".com/")[1]
	author := "<" + strings.Split(url, "/status")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	fxtwitterURL := strings.Replace(url, "twitter.com", "fxtwitter.com", 1)
	fxtwitterURL = strings.Replace(fxtwitterURL, "x.com", "fxtwitter.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Tweet](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, authorName, author, fxtwitterURL, mention)
	return fixedEmbedMessageContent
}

func fixRedditEmbed(m *discordgo.MessageCreate) string {
	url, fixedUrl := fixUrlEmbed(m)

	authorName := strings.Split(strings.Split(url, "/comments")[0], "r/")[1]
	author := "<" + strings.Split(url, "/comments")[0] + ">"
	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	vxredditURL := strings.Replace(url, "reddit.com", "vxreddit.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Reddit](%s) • [%s](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, authorName, author, vxredditURL, mention)
	return fixedEmbedMessageContent
}

func fixInstagramEmbed(m *discordgo.MessageCreate) string {
	url, fixedUrl := fixUrlEmbed(m)

	mention := fmt.Sprintf("<@%s>", m.Author.ID)
	ddinstagramURL := strings.Replace(url, "instagram.com", "ddinstagram.com", 1)

	fixedEmbedMessageContent := fmt.Sprintf("[Instagram](%s) • [Fix](%s) • Enviado por %s ", fixedUrl, ddinstagramURL, mention)
	return fixedEmbedMessageContent
}

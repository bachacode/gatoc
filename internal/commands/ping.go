package commands

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(ping.Metadata.Name, ping)
}

var ping bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "gatoping",
		Description: "Devuelve la latencia en MS",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		latency := s.HeartbeatLatency().Milliseconds()

		// Follow up with the actual latency
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("GatoPong! Latencia de %dms.", latency),
			},
		})
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return fmt.Errorf("Error responding to interaction: %v", err)
		}
		return nil
	},
}

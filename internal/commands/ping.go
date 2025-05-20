package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(ping.Metadata.Name, ping)
}

var ping SlashCommand = SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "gatoping",
		Description: "Devuelve la latencia en MS",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		latency := s.HeartbeatLatency().Milliseconds()

		// Follow up with the actual latency
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("GatoPong! Latencia de %dms.", latency),
			},
		})
		if err != nil {
			getInteractionFailedResponse(s, i, "")
			return fmt.Errorf("Error responding to interaction: %v", err)
		}
		return nil
	},
}

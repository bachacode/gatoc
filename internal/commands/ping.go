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
		Name: "gaplay",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "query",
				Description: "La URL de la canción",
				Required:    true,
			},
		},
		Description: "Reproduce una canción de Youtube.",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		latency := s.HeartbeatLatency().Milliseconds()

		// Follow up with the actual latency
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("GatoPong! Latencia de %dms.", latency),
			},
		})
		if err != nil {
			fmt.Println("Error sending follow-up message:", err)
		}
	},
}

package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Metadata *discordgo.ApplicationCommand
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Commands []SlashCommand = []SlashCommand{
	{
		Metadata: &discordgo.ApplicationCommand{
			Name:        "gatoping",
			Description: "Devuelve la latencia en MS",
		},
		Handler: PingCommandHandler,
	},
}

func PingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
}

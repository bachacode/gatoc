package events

import (
	"fmt"
	"os"
	"path/filepath"

	// Required for embed.FS
	"github.com/bwmarrin/discordgo"
)

func init() {
	register(guildMemberAdd)
}

var guildMemberAdd Event = Event{
	Name: "Guild Member Add / Join",
	Once: true,
	Handler: func(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
		channelID := os.Getenv("MAIN_CHANNEL_ID")
		env := os.Getenv("ENV")

		var emoji string
		if env == "development" {
			emoji = ":gatoc:1356083935460851878"
		} else {
			emoji = ":gatoc:1356257759850663976"
		}

		if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, "603340605774626871"); err != nil {
			fmt.Println("Failed to add default role:", err)
		}
		if r.Member.User.Username == "juanino" {
			if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, "603340605774626871"); err != nil {
				fmt.Println("Failed to add role to juanino:", err)
			}
		}

		// Path to the gif relative to the project root
		gifPath := filepath.Join("assets", "cat.gif")

		// Open the gif file
		file, err := os.Open(gifPath)
		if err != nil {
			fmt.Printf("Failed to open cat gif file: %v\n", err)
			return
		}
		defer file.Close()

		embed := discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "qlq <" + emoji + "> 🍷",
					Color:       0xFFFFFF,
					Description: r.Member.DisplayName() + " acaba de cometer el error mas grande de su vida entrando a esta tierra profana.",
					Image: &discordgo.MessageEmbedImage{
						URL: "attachment://cat.gif",
					},
				},
			},
			Files: []*discordgo.File{
				{
					Name:   "cat.gif",
					Reader: file,
				},
			},
		}
		_, err = s.ChannelMessageSendComplex(channelID, &embed)

		if err != nil {
			fmt.Printf("Failed to get main channel: %v\n", err)
		}
	},
}

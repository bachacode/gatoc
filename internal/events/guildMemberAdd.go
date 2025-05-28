package events

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bachacode/go-discord-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(guildMemberAdd)
}

var guildMemberAdd bot.Event = bot.Event{
	Name: "Guild Member Add / Join",
	Once: true,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
			channelID := ctx.MainChannelID
			emoji := ctx.WelcomeEmoji
			db := ctx.DB

			var wRoles []database.WelcomeRole
			result := db.Find(&wRoles)

			if result.Error != nil {
				fmt.Printf("Failed to get welcome roles: %v\n", result.Error)
				return
			}

			for _, wRole := range wRoles {
				if wRole.UserID == nil {
					if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, wRole.RoleID); err != nil {
						fmt.Println("Failed to add role to new member:", err)
					}
					continue
				}

				if *wRole.UserID == r.Member.User.ID {
					if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, wRole.RoleID); err != nil {
						fmt.Println("Failed to add role to new member:", err)
					}
				}
			}

			// if err := s.GuildMemberRoleAdd(r.GuildID, r.Member.User.ID, "603340605774626871"); err != nil {
			// 	fmt.Println("Failed to add default role:", err)
			// }

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
		}
	},
}

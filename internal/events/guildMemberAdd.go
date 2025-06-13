package events

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bachacode/go-discord-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(guildMemberAdd)
}

var guildMemberAdd bot.Event = bot.Event{
	Name: "Guild Member Add / Join",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
			channelID := ctx.MainChannelID
			emoji := ctx.WelcomeEmoji
			db := ctx.DB
			mention := "<@" + r.User.ID + ">"
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

			embed := discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "qlq " + emoji + " 🍷",
						Color:       0xFFFFFF,
						Description: mention + " acaba de cometer el error mas grande de su vida entrando a esta tierra profana.",
						Image: &discordgo.MessageEmbedImage{
							URL: "https://media.tenor.com/eH-RoS91Q1gAAAAC/cat.gif",
						},
					},
				},
			}
			_, err := s.ChannelMessageSendComplex(channelID, &embed)

			if err != nil {
				fmt.Printf("Failed to get main channel: %v\n", err)
			}
		}
	},
}

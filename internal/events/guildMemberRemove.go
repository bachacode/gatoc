package events

import (
	"fmt"
	"math/rand"

	"github.com/bachacode/gatoc/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterEvent(guildMemberRemove)
}

type LeaveMessage struct {
	embed    *discordgo.MessageEmbed
	filename string
}

var guildMemberRemove bot.Event = bot.Event{
	Name: "Guild Member Remove / Leave",
	Once: false,
	Handler: func(ctx *bot.BotContext) interface{} {
		return func(s *discordgo.Session, r *discordgo.GuildMemberRemove) {
			channelID := ctx.MainChannelID
			memberID := r.User.ID
			mention := "<@" + memberID + ">"
			messages := []LeaveMessage{
				{
					embed: &discordgo.MessageEmbed{
						Title:       "c lo acomodaron por las costillas <:sadcheems:869742943425151087>",
						Color:       0xFFFFFF,
						Description: mention + " no aguanto la pela.",
						Image: &discordgo.MessageEmbedImage{
							URL: "https://media.tenor.com/ww56Kix_vM8AAAAC/seloacomodoporlascostillas.gif",
						},
					},
					filename: "chavez.gif",
				},
				{
					embed: &discordgo.MessageEmbed{
						Title:       "c le fue la luz <:sadcheems:869742943425151087>",
						Color:       0xFFFFFF,
						Description: mention + " no aguanto la pela.",
						Image: &discordgo.MessageEmbedImage{
							URL: "https://media.tenor.com/vHMD9o7RmfYAAAAC/snake-salute.gif",
						},
					},
					filename: "snake.gif",
				},
			}

			randNumber := rand.Intn(len(messages))

			selectedMessage := messages[randNumber]

			embed := discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					selectedMessage.embed,
				},
			}
			_, err := s.ChannelMessageSendComplex(channelID, &embed)

			if err != nil {
				fmt.Printf("Failed to get main channel: %v\n", err)
			}
		}
	},
}

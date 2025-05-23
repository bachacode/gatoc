package commands

import (
	"strings"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(play.Metadata.Name, play)
}

var play bot.SlashCommand = bot.SlashCommand{
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
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		// get query string
		query, err := bot.GetInteractionOptionString("query", i)
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		// get voice state
		voiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "No estas conectado a un canal de voz!")
			return err
		}

		// defer reply
		var content string
		if err = bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		// get voice connection
		voiceConn, err := s.ChannelVoiceJoin(i.GuildID, voiceState.ChannelID, false, true)
		if err != nil {
			content = "Ha ocurrido un error al intentar unirse al canal."
			bot.EditDeferred(s, i, &content)
			return err
		}
		defer voiceConn.Disconnect()

		if !strings.Contains(query, "v=") {
			content = "La URL dada no tiene el formato correcto."
			bot.EditDeferred(s, i, &content)
			return err
		}

		videoID := strings.Split(query, "v=")[1]
		content = videoID
		bot.EditDeferred(s, i, &content)
		return nil
	},
}

package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(play.Metadata.Name, play)
}

var play SlashCommand = SlashCommand{
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
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) error {
		// get query string
		query, err := getInteractionOptionString("query", i)
		if err != nil {
			getInteractionFailedResponse(s, i, "")
			return err
		}

		// get voice state
		voiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
		if err != nil {
			getInteractionFailedResponse(s, i, "No estas conectado a un canal de voz!")
			return err
		}

		// defer reply
		var content string
		if err = deferReply(s, i); err != nil {
			getInteractionFailedResponse(s, i, "")
			return err
		}

		// get voice connection
		voiceConn, err := s.ChannelVoiceJoin(i.GuildID, voiceState.ChannelID, false, true)
		if err != nil {
			content = "Ha ocurrido un error al intentar unirse al canal."
			editDeferred(s, i, &content)
			return err
		}
		defer voiceConn.Disconnect()

		if !strings.Contains(query, "v=") {
			content = "La URL dada no tiene el formato correcto."
			editDeferred(s, i, &content)
			return err
		}

		videoID := strings.Split(query, "v=")[1]
		content = videoID
		editDeferred(s, i, &content)
		return nil
	},
}

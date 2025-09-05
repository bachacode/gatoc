package commands

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	subcommands "github.com/bachacode/gatoc/internal/commands/response-messages"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(responseMessage.Metadata.Name, responseMessage)
}

var responseMessage bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:                     "respuestas",
		Description:              "Gestiona las respuestas a mensajes especificos",
		DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.List.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "lista":
			return subcommands.List.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}

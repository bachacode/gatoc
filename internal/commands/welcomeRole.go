package commands

import (
	"fmt"

	"github.com/bachacode/gatoc/internal/bot"
	subcommands "github.com/bachacode/gatoc/internal/commands/welcome-roles"
	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.RegisterCommand(welcomeRole.Metadata.Name, welcomeRole)
}

var defaultMemberPermissions int64 = discordgo.PermissionManageServer

var welcomeRole bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:                     "rol-de-bienvenida",
		Description:              "Gestiona los roles otorgados a los nuevos miembros",
		DefaultMemberPermissions: &defaultMemberPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			subcommands.Add.Metadata,
			subcommands.List.Metadata,
			subcommands.Delete.Metadata,
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		switch options[0].Name {
		case "añadir":
			return subcommands.Add.Handler(s, i, ctx)
		case "lista":
			return subcommands.List.Handler(s, i, ctx)
		case "eliminar":
			return subcommands.Delete.Handler(s, i, ctx)
		default:
			bot.GetInteractionFailedResponse(s, i, "El subcomando llamado no existe.")
			return fmt.Errorf("Subcommand doesn't exist\n")
		}
	},
}

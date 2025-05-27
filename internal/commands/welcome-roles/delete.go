package welcomeroles

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bachacode/go-discord-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

var Delete bot.SlashSubcommand = bot.SlashSubcommand{
	Metadata: &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "eliminar",
		Description: "Elimina un rol de bienvenida",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole, // This is the key!
				Name:        "ID",
				Description: "Rol a eliminar",
				Required:    true,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {
		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options[0].Options {
			optionMap[opt.Name] = opt
		}

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		db := ctx.DB
		var roleID string
		var selectedRole *discordgo.Role
		var content string
		if roleOption, ok := optionMap["ID"]; ok {
			roleID = roleOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Roles != nil {
				selectedRole = i.ApplicationCommandData().Resolved.Roles[roleID]
			}
		} else {
			content := "Ha ocurrido un error para obtener el rol"
			bot.EditDeferred(s, i, &content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		if result := db.Delete(&database.WelcomeRole{}, roleID); result.Error != nil {
			content := "Ha ocurrido un error al eliminar el rol de bienvenida"
			bot.EditDeferred(s, i, &content)
			return fmt.Errorf("Error deleting welcome role: %s\n%v", selectedRole.Name, result.Error)
		}

		content = fmt.Sprintf("El rol `%s` ha sido eliminado de los roles de bienvenida", selectedRole.Name)
		bot.EditDeferred(s, i, &content)

		return nil
	},
}

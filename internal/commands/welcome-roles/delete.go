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
				Type:        discordgo.ApplicationCommandOptionString, // This is the key!
				Name:        "id",
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
		var wRoleID string
		var content string
		if roleOption, ok := optionMap["id"]; ok {
			wRoleID = roleOption.Value.(string)
		} else {
			content = "Ha ocurrido un error para obtener el rol."
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		if result := db.Delete(&database.WelcomeRole{}, wRoleID); result.Error != nil {
			content = "Ha ocurrido un error al eliminar el rol de bienvenida"
			bot.EditDeferred(s, i, content)
			return fmt.Errorf("Error deleting welcome role: %s\n%v", wRoleID, result.Error)
		} else if result.RowsAffected < 1 {
			content = "El ID introducido no existe en los roles de bienvenida"
			bot.EditDeferred(s, i, content)
			return nil
		}

		content = fmt.Sprintf("El rol de ID `%s` ha sido eliminado de los roles de bienvenida", wRoleID)
		bot.EditDeferred(s, i, content)

		return nil
	},
}

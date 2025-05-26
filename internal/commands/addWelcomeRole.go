package commands

import (
	"fmt"

	"github.com/bachacode/go-discord-bot/internal/bot"
	"github.com/bachacode/go-discord-bot/internal/database"
	"github.com/bwmarrin/discordgo"
)

// func init() {
// 	bot.RegisterCommand(addWelcomeRole.Metadata.Name, addWelcomeRole)
// }

var addWelcomeRole bot.SlashCommand = bot.SlashCommand{
	Metadata: &discordgo.ApplicationCommand{
		Name:        "rol-de-bienvenida",
		Description: "Añade un nuevo rol que se agregara a los miembros al entrar al servidor",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole, // This is the key!
				Name:        "role",
				Description: "Rol a añadir",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Usuario para añadir rol especifico (opcional)",
				Required:    false,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *bot.BotContext) error {

		// Get the permissions for the user in the channel where the command was sent
		perms, err := s.State.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			bot.GetInteractionFailedResponse(s, i, "No tienes permisos en este canal.")
			return err
		}

		// Check if the user has the Administrator permission
		if perms&discordgo.PermissionAdministrator != discordgo.PermissionAdministrator {
			bot.GetInteractionFailedResponse(s, i, "No tienes permisos para ejecutar este comando.")
			return err
		}

		options := i.ApplicationCommandData().Options

		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		if err := bot.DeferReply(s, i); err != nil {
			bot.GetInteractionFailedResponse(s, i, "")
			return err
		}

		db := ctx.DB
		guildID := ctx.GuildID
		var selectedRole *discordgo.Role
		var targetUser *discordgo.User
		var content string
		if roleOption, ok := optionMap["role"]; ok {
			roleID := roleOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Roles != nil {
				selectedRole = i.ApplicationCommandData().Resolved.Roles[roleID]
			}
		} else {
			content := "Ha ocurrido un error para obtener el rol"
			bot.EditDeferred(s, i, &content)
			return fmt.Errorf("Error responding to interaction\n")
		}

		// Get the "user" option (optional)
		if userOption, ok := optionMap["user"]; ok {
			userID := userOption.Value.(string)
			if i.ApplicationCommandData().Resolved != nil && i.ApplicationCommandData().Resolved.Users != nil {
				targetUser = i.ApplicationCommandData().Resolved.Users[userID]
			}
		}

		welcomeRole := database.WelcomeRole{
			GuildID: guildID,
			RoleID:  selectedRole.ID,
		}

		content = fmt.Sprintf("El rol `%s` será asignado a los nuevos miembros", selectedRole.Name)
		if targetUser != nil {
			content = fmt.Sprintf("El rol `%s` será asignado al usuario `%s`", selectedRole.Name, targetUser.Username)
			welcomeRole.UserID = &targetUser.ID
		}

		db.Create(&welcomeRole)

		bot.EditDeferred(s, i, &content)

		return nil
	},
}

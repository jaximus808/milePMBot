package functions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session, guildId string) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "project",
			Description: "Manage a project",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "op",
					Description:  "task operation",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "args",
					Description: "command arguments",
					Required:    false,
				},
			},
		},
		{
			Name:        "task",
			Description: "Manage a project's task",
			//make a better name for this shit
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "create a task",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "the new tasks's name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc",
							Description: "desc",
							Required:    true,
						},
					},
				}, {
					Name:        "assign",
					Description: "assign a task",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "assigned",
							Description: "the assigned user",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "taskref", // i should 100% make this auto complete
							Description: "the task ref",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "expectation",
							Description: "due date or story points",
							Required:    true,
						},
					},
				},
			},
		},
	}
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, cmd)
		if err != nil {
			log.Printf(" error %v", err)
		}
	}
}

func ClearCommands(s *discordgo.Session, guildId string) {
	cmds, err := s.ApplicationCommands(s.State.User.ID, guildId)

	if err != nil {
		log.Printf("Fialed to fetch existing commands: %v", err)
		return
	}
	for _, cmd := range cmds {
		s.ApplicationCommandDelete(s.State.User.ID, guildId, cmd.ID)

	}
}

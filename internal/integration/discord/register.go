package integration

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommands(s *discordgo.Session, guildId string, opGuildId string) {

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "project",
			Description: "Manage a project",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "start",
					Description: "Create a project (SERVER ADMIN ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "msname",
							Description: "the initial milestone's name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "msdate",
							Description: "due date",
							Required:    true,
						}, {
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc",
							Description: "description",
							Required:    true,
						},
					},
				},
				{
					Name:        "set",
					Description: "update a project setting (OWNER ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "setting",
							Description:  "the project setting to change",
							Required:     true,
							Autocomplete: true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "value",
							Description: "the new setting value",
							Required:    true,
						},
					},
				},
				{
					Name:        "role",
					Description: "update a project roles (ADMIN ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "op",
							Description:  "add a role/remove a role",
							Required:     true,
							Autocomplete: true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "user",
							Description: "user",
							Required:    true,
						},
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "role",
							Description:  "role: Required for adding a role",
							Required:     false,
							Autocomplete: true,
						},
					},
				},
				{
					Name:        "end",
					Description: "end a project (OWNER ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "move",
					Description: "end a project (OWNER ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "projectref",
							Description: "type the project ref for confirmation",
							Required:    true,
						},
					},
				},
				{
					Name:        "resume",
					Description: "end a project (OWNER ONLY)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "projectref",
							Description: "type the project ref for confirmation",
							Required:    true,
						},
					},
				},
				{
					Name:        "info",
					Description: "Get info on a project",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "milestone",
			Description: "Manage a project's milestones",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "create a milestone",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "msname",
							Description: "the initial milestone's name",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "msdate",
							Description: "due date",
							Required:    true,
						}, {
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc",
							Description: "description",
							Required:    true,
						},
					},
				},
				{
					Name:        "move",
					Description: "move to a milestone",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "direction",
							Description: "the direction to move the active milestone",
							Required:    true,
						},
					},
				}, {
					Name:        "delete",
					Description: "delete a milestone",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "msref",
							Description: "The milestone's ref",
							Required:    true,
						},
					},
				},
				{
					Name:        "map",
					Description: "create a milestone map",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
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
							Name:        "user",
							Description: "the assigned user",
							Required:    true,
						},
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "taskref", // i should 100% make this auto complete
							Description:  "the task ref (must start with milestone<id>/…)",
							Autocomplete: true,
							Required:     true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "expectation",
							Description: "due date or story points",
							Required:    true,
						},
					},
				}, {
					Name:        "complete",
					Description: "compte a task",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "taskref",
							Description:  "the complete taskref (must start with milestone<id>/…)",
							Autocomplete: true,
							Required:     true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc", // i should 100% make this auto complete
							Description: "progress description",
							Required:    true,
						},
					},
				}, {
					Name:        "approve",
					Description: "approve a task (LEADS+ only)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "taskref",
							Description:  "the completed taskref (must start with milestone<id>/…)",
							Autocomplete: true,
							Required:     true,
						},
					},
				}, {
					Name:        "reject",
					Description: "approve a task (LEADS+ only)",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "taskref",
							Description:  "the complete taskref",
							Autocomplete: true,
							Required:     true,
						}, {
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc",
							Description: "problems with the work",
							Required:    true,
						},
					},
				}, {
					Name:        "progress",
					Description: "adds progress to a task",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "taskref",
							Description:  "the complete taskref",
							Autocomplete: true,
							Required:     true,
						}, {
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "desc",
							Description: "problems with the work",
							Required:    true,
						},
					},
				}, {
					Name:        "list",
					Description: "list tasks and their status for a project, pass a user @ to get their assigned tasks",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "user",
							Description: "user to list tasks about",
						},
					},
				},
			},
		},
		{
			Name:        "help",
			Description: "manual for MilestonePM Bot",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "command",
					Description:  "create a milestone",
					Type:         discordgo.ApplicationCommandOptionString,
					Autocomplete: true,
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

	_, err := s.ApplicationCommandCreate(s.State.User.ID, opGuildId, &discordgo.ApplicationCommand{
		Name:        "deploywarn",
		Description: "warns all active projects that the bot will be down for 5 minutes for a deployment",
	})
	if err != nil {
		log.Printf(" error %v", err)
	}
	_, err = s.ApplicationCommandCreate(s.State.User.ID, opGuildId, &discordgo.ApplicationCommand{
		Name:        "hardstop",
		Description: "hardstops the serivce",
	})

	if err != nil {
		log.Printf(" error %v", err)
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

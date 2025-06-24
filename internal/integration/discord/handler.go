package integration

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/functions/general"
	"github.com/jaximus808/milePMBot/internal/functions/milestones"
	"github.com/jaximus808/milePMBot/internal/functions/projects"
	"github.com/jaximus808/milePMBot/internal/functions/tasks"
	"github.com/jaximus808/milePMBot/internal/util"
)

const CommandPrefix = "/"

var ActiveDBClient = util.SupaDB{}

// Autocomplete interaction handler
func autocompleteHandler(sess *discordgo.Session, interaction *discordgo.InteractionCreate) {

	cmd := interaction.ApplicationCommandData().Name
	switch cmd {
	case "task":
		handleTaskAutocomplete(sess, interaction)
	case "project":
		handleSettingAutocomplete(sess, interaction)
	case "help":
		handleHelpAutocomplete(sess, interaction)
	}
}

func handleHelpAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Find which option is focused
	// I need to implement a cache here to speed up performance
	for _, opt := range data.Options {
		if opt.Focused && opt.Name == "command" {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Project Commands",
							Value: "project",
						},
						{
							Name:  "Milestone Commands",
							Value: "milestone",
						},
						{
							Name:  "Task Commands",
							Value: "task",
						},
					},
				},
			})
			return
		}

	}
}
func handleSettingAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Find which option is focused
	// I need to implement a cache here to speed up performance
	for _, opt := range data.Options[0].Options {
		if opt.Focused && opt.Name == "setting" {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Output Channel [#channel]",
							Value: "output",
						},
						{
							Name:  "Project Description",
							Value: "description",
						},
						{
							Name:  "Sprint Message",
							Value: "message",
						},
						{
							Name:  "Toggle Sprint [y/n]",
							Value: "sprints",
						},
						{
							Name:  "Sprint Duration [# of weeks]",
							Value: "duration",
						},
						{
							Name:  "Toggle Weekly Pings [y/n]",
							Value: "pings",
						},
					},
				},
			})
			return
		}
		if opt.Focused && opt.Name == "op" {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Assign a Role",
							Value: "add",
						},
						{
							Name:  "Remove a Role",
							Value: "remove",
						},
					},
				},
			})
			return
		}
		if opt.Focused && opt.Name == "role" {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Lead: Ability to assign tasks to members and leads",
							Value: "lead",
						},
						{
							Name:  "Admin: Full power!!!!",
							Value: "admin",
						},
					},
				},
			})
			return
		}
	}
}

func handleTaskAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	// Find which option is focused
	// I need to implement a cache here to speed up performance
	for _, opt := range data.Options[0].Options {
		if opt.Focused && opt.Name == "taskref" {
			prefix := opt.StringValue()

			// first validate that the prefix is valid
			if !util.ValidTaskQuery(prefix) {
				return
			}

			subComamndName := data.Options[0].Name

			channel, errChannel := discord.DiscordSession.Channel(i.ChannelID)
			if errChannel != nil || channel == nil {
				return
			}

			if channel.ParentID == "" {
				return
			}

			activeProject, errActiveProject := ActiveDBClient.DBGetActiveProject(channel.GuildID, channel.ParentID)

			if errActiveProject != nil || activeProject == nil {
				log.Printf("No active project is running!")
				return
			}

			// if this is true, we're getting stories that are user has assigned
			var isAssigner = subComamndName == "reject" || subComamndName == "approve"

			// log.Printf("%s %s %t %d", i.Member.User.ID, prefix, isAssigner, *activeProject.ProjectID)

			var taskOptions *[]util.Task
			var taskOptionsError error

			// i need to restrict this to allow auto complete for leads
			if subComamndName == "assign" {
				taskOptions, taskOptionsError = ActiveDBClient.DBGetUnassignedTasks(i.Member.User.ID, prefix, *activeProject.ProjectID)
			} else {
				taskOptions, taskOptionsError = ActiveDBClient.DBGetTasksAndSpecifyDC(i.Member.User.ID, prefix, isAssigner, *activeProject.ProjectID, isAssigner, false)
			}

			if taskOptionsError != nil || taskOptions == nil {
				return
			}

			// // Fetch & filter your tasks/types from your DB:
			// matches := fetchTypeRefs(prefix) // implement to query: WHERE name ILIKE prefix||'%' LIMIT 25

			// // Build up to 25 choices
			var choices []*discordgo.ApplicationCommandOptionChoice
			for _, taskOption := range *taskOptions {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  *taskOption.TaskRef,
					Value: taskOption.TaskRef,
				})
			}

			// Send autocomplete response
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{Choices: choices},
			})
			return
		}
	}
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()

	var commandFunction func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport
	var exists bool

	baseCommand := data.Name

	var subCommand *discordgo.ApplicationCommandInteractionDataOption

	if len(data.Options) == 0 {
		commandFunction, exists = general.CommandMap[baseCommand]
		if !exists {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Didn't recongnize your command: " + baseCommand,
				},
			})
			return
		}
	} else {
		subCommand = data.Options[0]

		switch baseCommand {
		case "project":
			commandFunction, exists = projects.CommandMap[subCommand.Name]
		case "milestone":
			commandFunction, exists = milestones.CommandMap[subCommand.Name]
		case "task":
			commandFunction, exists = tasks.CommandMap[subCommand.Name]
		case "help":
			commandFunction, exists = general.CommandMap[baseCommand]
			// a hacky workaround to make help work without a subcommand
			subCommand = &discordgo.ApplicationCommandInteractionDataOption{
				Options: []*discordgo.ApplicationCommandInteractionDataOption{{
					Name:  "command",
					Type:  discordgo.ApplicationCommandOptionString,
					Value: subCommand.Value,
				}},
			}
		default:
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Didn't recongnize your command: " + baseCommand,
				},
			})
			return
		}
	}

	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Didn't recongnize your sub command: " + baseCommand + ": " + subCommand.Name,
			},
		})
		return
	}
	handleReport := commandFunction(i, subCommand, ActiveDBClient)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: handleReport.GetInfo(),
		},
	})
	if handleReport.NeedsOutput() {
		s.ChannelMessageSendEmbed(handleReport.GetOutputId(), handleReport.GetOutputMsg())
	}
}

// Command Handler
func MainHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commandHandler(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		autocompleteHandler(s, i)
	}

}

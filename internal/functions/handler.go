package functions

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/functions/milestones"
	"github.com/jaximus808/milePMBot/internal/functions/projects"
	"github.com/jaximus808/milePMBot/internal/functions/tasks"
	"github.com/jaximus808/milePMBot/internal/util"
)

const CommandPrefix = "/"

// Autocomplete interaction handler
func autocompleteHandler(sess *discordgo.Session, interaction *discordgo.InteractionCreate) {

	cmd := interaction.ApplicationCommandData().Name
	switch cmd {
	case "task":
		handleTaskAutocomplete(sess, interaction)
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
				log.Print("no options")
				return
			}

			subComamndName := data.Options[0].Name

			channel, errChannel := discord.DiscordSession.Channel(i.ChannelID)
			if errChannel != nil || channel == nil {
				log.Printf("failed to get channel")
				return
			}

			if channel.ParentID == "" {
				log.Printf("not in a channel category")
				return
			}

			activeProject, errActiveProject := util.DBGetActiveProject(channel.GuildID, channel.ParentID)

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
				taskOptions, taskOptionsError = util.DBGetUnassignedTasks(i.Member.User.ID, prefix, *activeProject.ProjectID)
			} else {
				taskOptions, taskOptionsError = util.DBGetSimillarTasksAssignedAndSpecifyDone(i.Member.User.ID, prefix, isAssigner, *activeProject.ProjectID, isAssigner)
			}

			if taskOptionsError != nil || taskOptions == nil {

				log.Print("no options")
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

	baseCommand := data.Name
	subCommand := data.Options[0]

	var commandFunction func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport
	var exists bool

	switch baseCommand {
	case "project":
		commandFunction, exists = projects.CommandMap[subCommand.Name]
	case "milestone":
		commandFunction, exists = milestones.CommandMap[subCommand.Name]
	case "task":
		commandFunction, exists = tasks.CommandMap[subCommand.Name]
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Didn't recongnize your command: " + baseCommand,
			},
		})
		return
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
	handleReport := commandFunction(i, subCommand)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: handleReport.GetInfo(),
		},
	})
	if handleReport.NeedsOutput() {
		s.ChannelMessageSend(handleReport.GetOutputId(), ">>> "+handleReport.GetOutputMsg())
	}
}

// Command Handler
func MainHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// if msg.Author.ID == sess.State.User.ID {
	// 	return
	// }

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		commandHandler(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		autocompleteHandler(s, i)
	}

}

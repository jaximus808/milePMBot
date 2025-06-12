package functions

import (
	"github.com/bwmarrin/discordgo"
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
	// options := i.ApplicationCommandData().Options
	// var focused *discordgo.ApplicationCommandInteractionDataOption

	// // Find the focused option
	// for _, opt := range options {
	// 	if opt.Focused {
	// 		focused = opt
	// 		break
	// 	}
	// }

	// if focused == nil || focused.Name != "op" {
	// 	log.Println("No focused option found or not 'name'")
	// 	return
	// }

	// query := strings.ToLower(focused.StringValue())

	// // Simulate dynamic task names (replace with DB query etc.)
	// allTasks := []string{"create", "assign", "done", "approve", "reject"}
	// var suggestions []*discordgo.ApplicationCommandOptionChoice

	// for _, task := range allTasks {
	// 	if strings.Contains(strings.ToLower(task), query) {
	// 		suggestions = append(suggestions, &discordgo.ApplicationCommandOptionChoice{
	// 			Name:  task,
	// 			Value: task,
	// 		})
	// 	}
	// 	if len(suggestions) >= 25 {
	// 		break
	// 	}
	// }

	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionApplicationCommandAutocompleteResult,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Choices: suggestions,
	// 	},
	// })
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	data := i.ApplicationCommandData()

	baseCommand := data.Name
	subCommand := data.Options[0]

	// for _, opt := range data.Options {
	// 	switch opt.Name {
	// 	case "op":
	// 		subCommand = opt.StringValue()
	// 	case "args":
	// 		args = strings.Fields(opt.StringValue())
	// 	}
	// }

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

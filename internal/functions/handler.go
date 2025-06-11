package functions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/functions/milestones"
	"github.com/jaximus808/milePMBot/internal/functions/projects"
	"github.com/jaximus808/milePMBot/internal/functions/tasks"
	"github.com/jaximus808/milePMBot/internal/util"
)

const CommandPrefix = "!"

func MainHandler(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}
	if strings.HasPrefix(msg.Content, CommandPrefix) {

		var commandFunction func(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport
		var exists bool

		args := strings.Split(msg.Content[1:], " ")

		if len(args) < 2 {
			sess.ChannelMessageSend(msg.ChannelID, "need more args")
			return
		}
		baseCommand := args[0]
		subCommand := args[1]
		switch baseCommand {
		case "project":
			commandFunction, exists = projects.CommandMap[subCommand]
		case "milestone":
			commandFunction, exists = milestones.CommandMap[subCommand]
		case "task":
			commandFunction, exists = tasks.CommandMap[subCommand]
		default:
			sess.ChannelMessageSend(msg.ChannelID, "Didn't recongize your command: "+args[0])
			return
		}

		if !exists {
			sess.ChannelMessageSend(msg.ChannelID, "Didn't recongize your command for "+args[0]+" : "+args[1])
			return
		}
		handleReport := commandFunction(msg, args[2:])
		sess.ChannelMessageSend(msg.ChannelID, handleReport.GetInfo())

		if handleReport.NeedsOutput() {
			sess.ChannelMessageSend(handleReport.GetOutputId(), ">>> "+handleReport.GetOutputMsg())
		}
	}
}

package functions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/functions/milestones"
	"github.com/jaximus808/milePMBot/internal/functions/projects"
)

const CommandPrefix = "!"

func MainHandler(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}
	if strings.HasPrefix(msg.Content, CommandPrefix) {

		args := strings.Split(msg.Content[1:], " ")
		switch args[0] {
		case "project":
			handleReport := projects.HandleCommand(msg, args[1:])
			sess.ChannelMessageSend(msg.ChannelID, handleReport.GetInfo())
			return
		case "milestone":
			handleReport := milestones.HandleCommand(msg, args[1:])
			sess.ChannelMessageSend(msg.ChannelID, handleReport.GetInfo())
			return
		default:
			sess.ChannelMessageSend(msg.ChannelID, "Didn't recongize your command: "+args[0])
		}

		// if true { // Replace with the report.success check
		// 	s.ChannelMessageSend(msg.ChannelID, "Success!") // Replace with report.info
		// } else {
		// 	s.ChannelMessageSend(m.ChannelID, "ERROR: Failed!") // Replace with report.info
		// }
	}
}

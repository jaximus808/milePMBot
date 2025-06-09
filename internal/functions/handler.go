package functions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const CommandPrefix = "!"

func MainHandler(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == sess.State.User.ID {
		return
	}
	if strings.HasPrefix(msg.Content, CommandPrefix) {

		args := strings.Split(msg.Content, " ")

		switch args[0] {
		case "project":

		}

		// if true { // Replace with the report.success check
		// 	s.ChannelMessageSend(msg.ChannelID, "Success!") // Replace with report.info
		// } else {
		// 	s.ChannelMessageSend(m.ChannelID, "ERROR: Failed!") // Replace with report.info
		// }
	}
}

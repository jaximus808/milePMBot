package general

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func HelpMsg(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	//todo: implement a job to delete these messages after 30 seconds to reduce clutter

	if args == nil {
		return util.CreateHandleReport(true, output.HELP_MSG)
	}

	command := util.GetOptionValue(args.Options, "command")

	switch strings.ToLower(command) {
	case "project":
		return util.CreateHandleReport(true, output.HELP_MSG_PROJECT)
	case "milestone":
		return util.CreateHandleReport(true, output.HELP_MSG_MILESTONE)
	case "task":
		return util.CreateHandleReport(true, output.HELP_MSG_TASK)
	}

	return util.CreateHandleReport(false, "❌ Didn't recongize the command to offer help on")
}

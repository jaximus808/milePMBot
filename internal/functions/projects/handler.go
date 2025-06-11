package projects

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func HandleCommand(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {
	if len(args) == 0 {
		return util.CreateHandleReport(false, "Needs arguments :(")
	}
	command := args[0]

	commandMethod, exists := commandMap[command]

	if !exists {
		return util.CreateHandleReport(false, "didn't recongize command: "+command)
	}

	//calls the command and removes the first element since we don't need it
	return commandMethod(msgInstance, args[1:])
}

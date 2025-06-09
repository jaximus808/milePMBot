package projects

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func HandleCommand(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {
	return util.CreateHandleReport(false, "Not yet implemnted")
}

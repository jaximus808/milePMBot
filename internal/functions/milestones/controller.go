package milestones

import (
	"github.com/bwmarrin/discordgo"
	milestones "github.com/jaximus808/milePMBot/internal/functions/milestones/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var commandMap = map[string](func(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport){
	"create": milestones.CreateMilestone,
}

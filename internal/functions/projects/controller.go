package projects

import (
	"github.com/bwmarrin/discordgo"
	projects "github.com/jaximus808/milePMBot/internal/functions/projects/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport){
	"start": projects.CreateProject,
}

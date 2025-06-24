package devops

import (
	"github.com/bwmarrin/discordgo"
	devops "github.com/jaximus808/milePMBot/internal/functions/devops/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport){
	"warn":      devops.WarnDeployments,
	"forcestop": devops.HardStop,
}

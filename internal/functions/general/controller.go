package general

import (
	"github.com/bwmarrin/discordgo"
	general "github.com/jaximus808/milePMBot/internal/functions/general/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport){
	"help": general.HelpMsg,
}

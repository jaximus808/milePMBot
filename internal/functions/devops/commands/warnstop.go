package devops

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func WarnDeployments(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {
	// fuck i gotta do some bs here
	return util.CreateHandleReport(false, "WIP")
}

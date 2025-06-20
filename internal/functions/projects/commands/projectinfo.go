package projects

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ProjectInfo(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	return util.CreateHandleReport(true, fmt.Sprintf("# Project: %s\n## Desc: %s", *currentProject.ProjectRef, *currentProject.Desc))
}

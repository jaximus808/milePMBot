package projects

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ProjectInfo(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	// getting info
	milestone, milestoneErr := DB.DBGetMilestoneWithId(*currentProject.CurrentMID)

	if milestoneErr != nil || milestone == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReport(
		true, fmt.Sprintf(">>> # Project Ref: %s\n**Desc:** %s\nCurrent Milestone: %s\nDue: %s", *currentProject.ProjectRef, *currentProject.Desc, *milestone.DisplayName, milestone.DueDate.Format("January 2, 2006")))
}

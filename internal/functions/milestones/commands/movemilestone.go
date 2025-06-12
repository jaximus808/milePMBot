package milestones

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func MoveMilestone(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	currentMilestone, errorCurrentMilestone := util.DBGetMilestoneWithId(*currentProject.CurrentMID)
	if errorCurrentMilestone != nil {
		return util.CreateHandleReport(false, errorCurrentMilestone.Error())
	}

	if currentMilestone == nil {
		return util.CreateHandleReport(false, "failed to get active milestone")
	}
	direction := util.GetOptionValue(args.Options, "direction")

	var newMilestone *util.Milestone
	var errorNewMilestone error

	var displayDirection string

	if direction == "next" {
		//need to insert current project
		newMilestone, errorNewMilestone = util.DBGetNextMilestone(currentProject.ID, currentMilestone)
		displayDirection = "forward"
	} else if direction == "prev" {
		newMilestone, errorNewMilestone = util.DBGetPrevMilestone(currentProject.ID, currentMilestone)
		displayDirection = "back"
	}

	if errorNewMilestone != nil || newMilestone == nil {
		return util.CreateHandleReport(false, "Could not get "+direction+" milestone")
	}

	newProject, newProjectError := util.DBUpdateCurrentMilestone(currentProject.ID, newMilestone.ID)
	if newProjectError != nil {
		return util.CreateHandleReport(false, "Could not update the project to its new milestone")
	}

	if newProject == nil {
		return util.CreateHandleReport(false, "failed to update new milestone")
	}
	// will need to print tasks after
	return util.CreateHandleReportAndOutput(
		true,
		"successfully changed to milestone: "+(*newMilestone.DisplayName),
		fmt.Sprintf("**Milestone Update!** \nThe project's current milestone has now moved %s to %s", displayDirection, (*newMilestone.DisplayName)),
		*currentProject.OutputChannel,
	)
}

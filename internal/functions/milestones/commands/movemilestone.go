package milestones

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func MoveMilestone(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {

	if len(args) != 1 {
		return util.CreateHandleReport(false, "expected [next/prev]")
	}

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
	direction := args[0]

	var newMilestone *util.Milestone
	var errorNewMilestone error

	if direction == "next" {
		//need to insert current project
		newMilestone, errorNewMilestone = util.DBGetNextMilestone(currentProject.ID, currentMilestone)
	} else if direction == "prev" {
		newMilestone, errorNewMilestone = util.DBGetPrevMilestone(currentMilestone.ID, currentMilestone)
	}

	if errorNewMilestone != nil || newMilestone == nil {
		return util.CreateHandleReport(false, "Could not get "+direction+" project")
	}

	newProject, newProjectError := util.DBUpdateCurrentMilestone(currentMilestone.ID, newMilestone.ID)
	if newProjectError != nil {
		return util.CreateHandleReport(false, newProjectError.Error())
	}
	if newProject == nil {
		return util.CreateHandleReport(false, "failed to update new milestone")
	}
	return util.CreateHandleReport(true, "successfully changed to milestone: "+(*newMilestone.DisplayName))
}

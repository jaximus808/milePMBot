package tasks

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func RejectTask(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {

	if len(args) < 2 {
		return util.CreateHandleReport(false, "Expected [task_ref]")
	}

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	taskRef := args[0]
	desc := strings.Join(args[1:], " ")

	//pretty much gonna make a new progress

	currentTask, currentTaskError := util.DBGetTask(currentProject.ID, taskRef)

	if currentTask == nil || currentTaskError != nil {
		return util.CreateHandleReport(false, "failed to get ticket, check your task_ref!")
	}

	newProgress, newProgressError := util.DBCreateProgress(currentTask.ID, fmt.Sprintf("Not Approved: %s", desc), false)
	if newProgress == nil || newProgressError != nil {
		return util.CreateHandleReport(false, "something went wrong on our end :/")
	}
	return util.CreateHandleReportAndOutput(
		true,
		"We'll mark this as not approved and notify the assigned person",
		fmt.Sprintf("Task: %s **Not Approved**\nTask Ref: %s\nReason: %s\n<@%s> Please review and remedy these changes", *currentTask.TaskName, *currentTask.TaskRef, desc, *currentTask.AssignedID),
		*currentProject.OutputChannel,
	)
}

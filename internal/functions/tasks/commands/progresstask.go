package tasks

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ProgressTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	currentTask, currentTaskError := util.DBGetTask(currentProject.ID, taskRef)

	if currentTask == nil || currentTaskError != nil {
		return util.CreateHandleReport(false, "failed to get task, check your task_ref!")
	}

	newProgress, newProgressError := util.DBCreateProgress(currentTask.ID, fmt.Sprintf("Progress Update (%s): %s", time.Now().Format(time.RFC1123), desc), false)
	if newProgress == nil || newProgressError != nil {
		return util.CreateHandleReport(false, "something went wrong on our end :/")
	}
	return util.CreateHandleReportAndOutput(
		true,
		"Got it! Updated progress and letting your assigner know",
		fmt.Sprintf("Task: %s **Not Approved**\nTask Ref: %s\nReason: %s\n<@%s> Please review and remedy these changes", *currentTask.TaskName, *currentTask.TaskRef, desc, *currentTask.AssignedID),
		*currentProject.OutputChannel,
	)

}

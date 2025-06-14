package tasks

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func RejectTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	//pretty much gonna make a new progress

	currentTask, currentTaskError := util.DBGetTask(currentProject.ID, taskRef)

	if currentTask == nil || currentTaskError != nil {
		return util.CreateHandleReport(false, "failed to get task, check your task_ref!")
	}
	updatedTask, updatedTaskError := util.DBUpdateTaskRecentProgress(currentTask.ID, false)
	if updatedTaskError != nil || updatedTask == nil {
		return util.CreateHandleReport(false, "Failed to update task correctly")
	}

	newProgress, newProgressError := util.DBCreateProgress(currentTask.ID, fmt.Sprintf("Not Approved: %s", desc), false)
	if newProgress == nil || newProgressError != nil {
		return util.CreateHandleReport(false, "something went wrong on our end :/")
	}
	return util.CreateHandleReportAndOutput(
		true,
		"We'll mark this as not approved and notify the assigned person",
		&discordgo.MessageEmbed{
			Title:       "❌ Task Not Approved",
			Description: fmt.Sprintf("<@%s> did not approve the task. Please review and make necessary changes.", msgInstance.Member.User.ID),
			Color:       0xE74C3C, // Red
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *updatedTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *updatedTask.TaskRef, Inline: false},
				{Name: "Feedback", Value: desc, Inline: false}, // optional feedback
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		*currentProject.OutputChannel,
	)
}

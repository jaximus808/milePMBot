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
		&discordgo.MessageEmbed{
			Title:       "📈 Progress Update",
			Description: fmt.Sprintf("<@%s> added progress to a task.", *currentTask.AssignedID),
			Color:       0xF1C40F, // Yellow
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *currentTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *currentTask.TaskRef, Inline: false},
				{Name: "Progress Note", Value: desc, Inline: false}, // Assume `progressNote` is a string
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		*currentProject.OutputChannel,
	)

}

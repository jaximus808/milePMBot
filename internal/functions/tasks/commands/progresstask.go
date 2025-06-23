package tasks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ProgressTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)
	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	currentTask, currentTaskError := DB.DBGetTask(currentProject.ID, taskRef)

	if currentTask == nil || currentTaskError != nil {
		return util.CreateHandleReport(false, "failed to get task, check your task_ref!")
	}

	newProgress, newProgressError := DB.DBCreateProgress(currentTask.ID, fmt.Sprintf("Progress Update (%s): %s", time.Now().Format(time.RFC1123), desc), false)
	if newProgress == nil || newProgressError != nil {
		return util.CreateHandleReport(false, "something went wrong on our end :/")
	}
	return util.CreateHandleReportAndOutput(
		true,
		output.SUCCESS_PROGRESS_ADDED,
		&discordgo.MessageEmbed{
			Title:       "📈 Progress Update",
			Description: fmt.Sprintf("<@%d> added progress to a task.", *currentTask.AssignedID),
			Color:       0xF1C40F, // Yellow
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *currentTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *currentTask.TaskRef, Inline: false},
				{Name: "Progress Note", Value: desc, Inline: false}, // Assume `progressNote` is a string
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		strconv.Itoa(*currentProject.OutputChannel),
	)

}

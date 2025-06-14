package tasks

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ApproveTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	taskRef := util.GetOptionValue(args.Options, "taskref")

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	currentTime := time.Now()

	updatedTask, errorUpdatedTask := util.DBTaskMarkComplete(currentProject.ID, taskRef, &currentTime)

	if updatedTask == nil || errorUpdatedTask != nil {
		return util.CreateHandleReport(false, "Failed to mark task as done :(")
	}

	return util.CreateHandleReportAndOutput(
		true,
		"Yay! Task is now marked as approved :smile:",
		&discordgo.MessageEmbed{
			Title:       "🎉 Task Approved",
			Description: fmt.Sprintf("<@%s> approved the completion of a task.", msgInstance.Member.User.ID),
			Color:       0x9B59B6, // Purple
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *updatedTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *updatedTask.TaskRef, Inline: false},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		*currentProject.OutputChannel,
	)
}

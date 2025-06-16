package tasks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ApproveTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	taskRef := util.GetOptionValue(args.Options, "taskref")

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	currentTime := time.Now()

	updatedTask, errorUpdatedTask := util.DBTaskMarkComplete(currentProject.ID, taskRef, &currentTime)

	if updatedTask == nil || errorUpdatedTask != nil {
		return util.CreateHandleReport(false, output.FAIL_TASK_DNE)
	}

	return util.CreateHandleReportAndOutput(
		true,
		output.SUCCESS_APPROVING_TASK,
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
		strconv.Itoa(*currentProject.OutputChannel),
	)
}

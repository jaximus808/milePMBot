package tasks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CompleteTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	currentTask, errorCurrentTask := DB.DBGetTask(currentProject.ID, taskRef)

	if errorCurrentTask != nil || currentTask == nil {
		return util.CreateHandleReport(false, output.FAIL_TASK_DNE)
	}

	userId, err := strconv.Atoi(msgInstance.Member.User.ID)

	if err != nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	// now check if this task is even assigned to the user
	if *currentTask.AssignedID != userId {
		return util.CreateHandleReport(false, output.ERROR_NOT_YOUR_TASK)
	}

	//now we need to make a progress row then ask for review
	newProgress, newProgressError := DB.DBCreateProgress(currentTask.ID, desc, true)

	if newProgressError != nil || newProgress == nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	updatedTask, updatedTaskError := DB.DBUpdateTaskRecentProgress(currentTask.ID, true)

	if updatedTaskError != nil || updatedTask == nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	//should make the option a helper function lol
	return util.CreateHandleReportAndOutput(
		true,
		output.SUCCESS_COMPLETE_TASK,
		&discordgo.MessageEmbed{
			Title:       "✅ Task Completed",
			Description: fmt.Sprintf("<@%d> has completed a task.", *updatedTask.AssignedID),
			Color:       0x2ECC71, // Green
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *updatedTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *updatedTask.TaskRef, Inline: false},
				{Name: "Notes", Value: desc, Inline: false},
				{Name: "For Assigner", Value: fmt.Sprintf("<%d> Review this work and approve/disapprove with /task [approve/disapprove]", *updatedTask.AssignerID), Inline: false},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		strconv.Itoa(*currentProject.OutputChannel),
	)
}

package tasks

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CompleteTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return util.CreateHandleReport(false, "no project exists here :(")
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "no project exists here :(")
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	currentTask, errorCurrentTask := util.DBGetTask(currentProject.ID, taskRef)

	if errorCurrentTask != nil || currentTask == nil {
		return util.CreateHandleReport(false, "Invalid task ref")
	}

	// now check if this task is even assigned to the user
	if *currentTask.AssignedID != msgInstance.Member.User.ID {
		return util.CreateHandleReport(false, "This task isn't assigned to you")
	}

	//now we need to make a progress row then ask for review
	newProgress, newProgressError := util.DBCreateProgress(currentTask.ID, desc, true)

	if newProgressError != nil || newProgress == nil {
		return util.CreateHandleReport(false, "Couldn't make a progress report :()")
	}
	updatedTask, updatedTaskError := util.DBUpdateTaskRecentProgress(currentTask.ID, true)

	if updatedTaskError != nil || updatedTask == nil {
		return util.CreateHandleReport(false, "Failed to update task correctly")
	}

	//should make the option a helper function lol
	return util.CreateHandleReportAndOutput(
		true,
		"Marked as completed and sent to your assigner for review! ",
		&discordgo.MessageEmbed{
			Title:       "✅ Task Completed",
			Description: fmt.Sprintf("<@%s> has completed a task.", *updatedTask.AssignedID),
			Color:       0x2ECC71, // Green
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *updatedTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *updatedTask.TaskRef, Inline: false},
				{Name: "Notes", Value: desc, Inline: false},
				{Name: "For Assigner", Value: fmt.Sprintf("<%s> Review this work and approve/disapprove with /task [approve/disapprove]", *updatedTask.AssignerID), Inline: false},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		*currentProject.OutputChannel,
	)
}

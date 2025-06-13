package tasks

import (
	"fmt"
	"log"

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

	log.Printf("got task ref %s", taskRef)

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
		fmt.Sprintf("Task **%s** Completed!\nTask Ref: %s\n<@%s> Please review and approve/reject with !task [approve/reject] %s\nDesc: %s", *currentTask.TaskName, *currentTask.TaskRef, *currentTask.AssignerID, *currentTask.TaskRef, desc),
		*currentProject.OutputChannel,
	)
}

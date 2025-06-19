package tasks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func RejectTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)
	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	taskRef := util.GetOptionValue(args.Options, "taskref")
	desc := util.GetOptionValue(args.Options, "desc")

	//pretty much gonna make a new progress

	currentRole, currentRoleError := util.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if currentRoleError != nil || currentRole == nil {
		return util.CreateHandleReport(false, "❌ You lack valid permissions to approve tasks")
	}

	currentTask, currentTaskError := util.DBGetTask(currentProject.ID, taskRef)

	if currentTaskError != nil || currentTask == nil {
		return util.CreateHandleReport(false, fmt.Sprintf("❌ Could not find a task with the given task ref: %s", taskRef))
	}

	if currentTask.AssignedID == nil {
		return util.CreateHandleReport(false, "❌ This project has not been assigned yet")
	}

	if currentRole.RoleLevel == int(util.LeadRole) && strconv.Itoa(*currentTask.AssignedID) != msgInstance.Member.User.ID {
		return util.CreateHandleReport(false, "❌ As a lead you can only reject tasks you've assigned!")
	}
	updatedTask, updatedTaskError := util.DBUpdateTaskRecentProgress(currentTask.ID, false)
	if updatedTaskError != nil || updatedTask == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	newProgress, newProgressError := util.DBCreateProgress(currentTask.ID, fmt.Sprintf("Not Approved: %s", desc), false)
	if newProgress == nil || newProgressError != nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReportAndOutput(
		true,
		output.SUCCESS_REJECT,
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
		strconv.Itoa(*currentProject.OutputChannel),
	)
}

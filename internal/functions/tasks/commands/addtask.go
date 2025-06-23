package tasks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func AddTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	// ADD THE UNDERSCORE STUFF

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)
	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	userRole, errorRole := DB.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if errorRole != nil || userRole == nil {
		return util.CreateHandleReport(false, output.FAIL_PERMS)
	}

	if userRole.RoleLevel < int(util.LeadRole) {
		return util.CreateHandleReport(false, output.FAIL_PERMS)
	}

	// user has permisisons now

	// need to check if the same task ref already exists
	// need to also check if a milestone on the same date exists

	// Tbh i think it would be worth making a prompt feature, where the
	// bot can ask which task would u like to assign, and you can reply with the number
	// but then what if there are too many? FUCK, for now a shitty ref should work

	// ref can be milestone{milestoneId}/task_ref_name

	taskName := util.GetOptionValue(args.Options, "name")
	taskDesc := util.GetOptionValue(args.Options, "desc")

	newTask, taskError := DB.DBCreateTasks(currentProject.ID, taskName, taskDesc, *currentProject.CurrentMID)

	if taskError != nil || newTask == nil {
		return util.CreateHandleReport(false, output.FAIL_CREATE_TASK)
	}

	return util.CreateHandleReportAndOutput(
		true,
		fmt.Sprintf("Created task with task_ref: %s ",
			*newTask.TaskRef),
		&discordgo.MessageEmbed{
			Title:       "🆕 Task Created",
			Description: "A new task has been added to the project.",
			Color:       0x3498DB, // Blue
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: taskName, Inline: false},
				{Name: "Task Ref", Value: *newTask.TaskRef, Inline: false},
				{Name: "Description", Value: taskDesc, Inline: false},
				{Name: "Created By", Value: fmt.Sprintf("<@%s>", msgInstance.Member.User.ID), Inline: false},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		strconv.Itoa(*currentProject.OutputChannel),
	)
}

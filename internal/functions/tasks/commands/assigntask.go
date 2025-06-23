package tasks

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func AssignTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	assignedID := util.GetOptionValue(args.Options, "user")

	taskRef := util.GetOptionValue(args.Options, "taskref")
	taskExpectations := util.GetOptionValue(args.Options, "expectation")

	// if len(args) != 3 {
	// 	return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
	// }

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}
	re := regexp.MustCompile(`<@!?(\d+)>`)
	match := re.FindStringSubmatch(assignedID)
	if len(match) != 2 {
		return util.CreateHandleReport(false, "You need to be @ a user!")
	}
	assignedUserId := match[1]

	dueDate, dateError := time.Parse("01/02/2006", taskExpectations)

	var moreDetails string
	var assignedTask *util.Task
	var assignedError error
	// now do assignments

	//TODO: Make a better output prompt or just data dump
	// Need also a check to make sure the task isn't already done

	userId, idError := strconv.Atoi(msgInstance.Member.User.ID)
	assignedUserIdInt, idAssignError := strconv.Atoi(assignedUserId)
	if idError != nil || idAssignError != nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	//fix this
	if dateError != nil {
		storyPoint, errStoryPoint := strconv.Atoi(taskExpectations)
		if errStoryPoint != nil {
			return util.CreateHandleReport(false, output.ERROR_ARGS_ASSIGN)
		}

		assignedTask, assignedError = DB.DBAssignTasksStoryPoints(currentProject.ID, taskRef, userId, assignedUserIdInt, storyPoint)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, output.FAIL_TASK_DNE)
		}
		moreDetails = fmt.Sprintf("Story Points: %d", storyPoint)

	} else {
		assignedTask, assignedError = DB.DBAssignTasksDueDate(currentProject.ID, taskRef, userId, assignedUserIdInt, &dueDate)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, output.FAIL_TASK_DNE)
		}
		moreDetails = fmt.Sprintf("Due Date: %s", dueDate)
	}

	return util.CreateHandleReportAndOutput(
		true,
		fmt.Sprintf("Task %s assigned successfully", taskRef),
		&discordgo.MessageEmbed{
			Title:       "📌 Task Assigned",
			Description: "A task has been assigned to a team member.",
			Color:       0xE67E22, // Orange
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Task", Value: *assignedTask.TaskName, Inline: false},
				{Name: "Task Ref", Value: *assignedTask.TaskRef, Inline: false},
				{Name: "Assigned To", Value: fmt.Sprintf("<@%s>", assignedUserId), Inline: false},
				{Name: "By", Value: fmt.Sprintf("<@%s>", msgInstance.Member.User.ID), Inline: false},
				{Name: "Details", Value: moreDetails, Inline: false},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
		strconv.Itoa(*currentProject.OutputChannel),
	)

}

package tasks

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func AssignTask(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	assignedID := util.GetOptionValue(args.Options, "assigned")

	taskRef := util.GetOptionValue(args.Options, "taskref")
	taskExpectations := util.GetOptionValue(args.Options, "expectation")

	// if len(args) != 3 {
	// 	return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
	// }

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

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

	// now do assignments

	//TODO: Make a better output prompt or just data dump
	// Need also a check to make sure the task isn't already done
	if dateError != nil {
		storyPoint, errStoryPoint := strconv.Atoi(taskExpectations)
		if errStoryPoint != nil {
			return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
		}

		assignedTask, assignedError := util.DBAssignTasksStoryPoints(currentProject.ID, taskRef, msgInstance.Member.User.ID, assignedUserId, storyPoint)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, "Invalid task_ref, are you sure this exists for the current milestone?")
		}
		moreDetails = fmt.Sprintf("Story Points: %d", storyPoint)

	} else {
		assignedTask, assignedError := util.DBAssignTasksDueDate(currentProject.ID, taskRef, msgInstance.Member.User.ID, assignedUserId, &dueDate)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, "Invalid task_ref, are you sure this exists for the current milestone?")
		}
		moreDetails = fmt.Sprintf("Due Date: %s", dueDate)
	}

	return util.CreateHandleReportAndOutput(
		true,
		fmt.Sprintf("Task %s assigned successfully", taskRef),
		fmt.Sprintf("**Assigned** task_ref **%s** \n**To user:** <@%s>\n%s", taskRef, assignedUserId, moreDetails),
		*currentProject.OutputChannel,
	)

}

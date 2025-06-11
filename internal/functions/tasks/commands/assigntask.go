package tasks

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func AssignTask(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {
	if len(args) != 3 {
		return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
	}

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	re := regexp.MustCompile(`<@!?(\d+)>`)
	foundMatches := re.FindAllString(args[0], 1)

	if foundMatches == nil || len(foundMatches) != 0 {
		return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
	}

	assignedUserId := foundMatches[0]

	taskRef := args[1]

	dueDate, dateError := time.Parse("01/02/2006", args[2])

	// now do assignments

	//TODO: Make a better output prompt or just data dump
	// Need also a check to make sure the task isn't already done
	if dateError != nil {
		storyPoint, errStoryPoint := strconv.Atoi(args[2])
		if errStoryPoint != nil {
			return util.CreateHandleReport(false, "Expecting 3 args [@assign] [task_ref] [due_date or story points]")
		}

		assignedTask, assignedError := util.DBAssignTasksStoryPoints(currentProject.ID, taskRef, msgInstance.Author.ID, assignedUserId, storyPoint)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, "Invalid task_ref, are you sure this exists for the current milestone?")
		}

	} else {
		assignedTask, assignedError := util.DBAssignTasksDueDate(currentProject.ID, taskRef, msgInstance.Author.ID, assignedUserId, &dueDate)

		if assignedError != nil || assignedTask == nil {
			return util.CreateHandleReport(false, "Invalid task_ref, are you sure this exists for the current milestone?")
		}
	}

	return util.CreateHandleReport(true, fmt.Sprintf("Assigned task_ref %s to user: <@%s> successfully", taskRef, assignedUserId))

}

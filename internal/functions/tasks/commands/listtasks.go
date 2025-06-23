package tasks

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ListTasks(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	userRef := util.GetOptionValue(args.Options, "user")

	if userRef == "" {
		return getTaskForMilestone(msgInstance, DB)
	} else {
		re := regexp.MustCompile(`<@!?(\d+)>`)
		match := re.FindStringSubmatch(userRef)
		if len(match) != 2 {
			return util.CreateHandleReport(false, "You need to @ a user!")
		}
		userId := match[1]
		return getTaskForMilestoneForUser(msgInstance, userId, DB)
	}
}

func getTaskForMilestone(i *discordgo.InteractionCreate, DB util.DBClient) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(i, DB)

	if errorHandle != nil {
		return util.CreateHandleReport(false, "no project exists here :(")
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "no project exists here :(")
	}

	emeddedMessage := &discordgo.MessageEmbed{
		Title:       "📋 Task List",
		Description: "For: current milestone",
		Color:       0x4A90E2, // Green
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	tasks, tasksError := DB.DBGetTasksForMilestone(currentProject.ID, *currentProject.CurrentMID)

	taskReport := util.ParseTaskList(tasks, i.GuildID)

	doneField := &discordgo.MessageEmbedField{
		Name: "✅ Done",
	}
	inReviewField := &discordgo.MessageEmbedField{
		Name: "🔍 In Review",
	}
	inProgressField := &discordgo.MessageEmbedField{
		Name: "🛠️ In Progress",
	}
	unassignedField := &discordgo.MessageEmbedField{
		Name: "🗂 Backlog",
	}

	if len(taskReport.Done) > 0 {
		doneField.Value = strings.Join(taskReport.Done, "\n")
	} else {
		doneField.Value = "No Done Tasks"
	}

	if len(taskReport.InReview) > 0 {
		inReviewField.Value = strings.Join(taskReport.InReview, "\n")
	} else {
		inReviewField.Value = "No In Review Tasks"
	}

	if len(taskReport.InProgress) > 0 {
		inProgressField.Value = strings.Join(taskReport.InProgress, "\n")
	} else {
		inProgressField.Value = "No In Progress Tasks"
	}

	if len(taskReport.Unnassigned) > 0 {
		unassignedField.Value = strings.Join(taskReport.Unnassigned, "\n")
	} else {
		unassignedField.Value = "No Backlog Tasks"
	}

	if tasksError != nil || tasks == nil {
		return util.CreateHandleReport(false, "There are no tasks for the milestone!")
	}

	emeddedMessage.Fields = []*discordgo.MessageEmbedField{
		doneField,
		inReviewField,
		inProgressField,
		unassignedField,
	}
	return util.CreateHandleReportAndOutput(
		true,
		"List Made!",
		emeddedMessage,
		i.ChannelID,
	)
}

func getTaskForMilestoneForUser(i *discordgo.InteractionCreate, userId string, DB util.DBClient) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(i, DB)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	emeddedMessage := &discordgo.MessageEmbed{
		Title:       "📋 Task List",
		Description: fmt.Sprintf("For: %s", util.GetUserGuildNickname(i.GuildID, userId)),
		Color:       0x4A90E2, // Green
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	tasks, tasksError := DB.DBGetTasksForMilestoneAndAssignedUser(currentProject.ID, *currentProject.CurrentMID, userId)
	if tasksError != nil || tasks == nil {
		return util.CreateHandleReport(false, output.ERROR_USER_NO_TASK)
	}

	taskReport := util.ParseTaskList(tasks, i.GuildID)

	doneField := &discordgo.MessageEmbedField{
		Name: "✅ Done",
	}

	inReviewField := &discordgo.MessageEmbedField{
		Name: "🔍 In Review",
	}

	inProgressField := &discordgo.MessageEmbedField{
		Name: "🛠️ In Progress",
	}

	if len(taskReport.Done) > 0 {
		doneField.Value = strings.Join(taskReport.Done, "\n")
	} else {
		doneField.Value = "No Done Tasks"
	}

	if len(taskReport.InReview) > 0 {
		inReviewField.Value = strings.Join(taskReport.InReview, "\n")
	} else {
		inReviewField.Value = "No In Review Tasks"
	}

	if len(taskReport.InProgress) > 0 {
		inProgressField.Value = strings.Join(taskReport.InProgress, "\n")
	} else {
		inProgressField.Value = "No In Progress Tasks"
	}

	emeddedMessage.Fields = []*discordgo.MessageEmbedField{
		doneField,
		inReviewField,
		inProgressField,
	}

	return util.CreateHandleReportAndOutput(
		true,
		output.SUCCESS_TASK_LIST,
		emeddedMessage,
		i.ChannelID,
	)
}

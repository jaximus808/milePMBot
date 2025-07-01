package util

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
)

// TODO: Need to modify the error stuff, maybe also split this up into different files

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
func ptrInt(i int) *int          { return &i }

// milestones is assumed to be inorder
func ParseMilestoneList(milestones *[]Milestone, currentMid int) *MilestoneReport {

	milestoneMap := &MilestoneReport{}

	upcomingParsing := true
	// get the id of te curent MID

	for _, milestone := range *milestones {
		if milestone.ID == currentMid {
			milestoneMap.CurrentMilestone = fmt.Sprintf("> Milestone: %s\n> Description: %s", *milestone.DisplayName, *milestone.Description)
			upcomingParsing = false

		} else if upcomingParsing {
			milestoneMap.Upcoming = append(
				milestoneMap.Upcoming,
				fmt.Sprintf("> Milestone: %s\n> Description: %s", *milestone.DisplayName, *milestone.Description),
			)
		} else {
			if upcomingParsing {
				milestoneMap.Previous = append(
					milestoneMap.Previous,
					fmt.Sprintf("> Milestone: %s\n> Description: %s", *milestone.DisplayName, *milestone.Description),
				)
			}
		}

	}

	return milestoneMap

}

func ParseTaskListWeeklyWithPing(tasks *[]Task) *TaskReport {

	taskReport := &TaskReport{}

	for _, task := range *tasks {

		if *task.Completed { // task is marked compelete by the user and ready for review

			taskReport.InReview = append(
				taskReport.InReview,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Completed By: <@%d>\n> Waiting for Review By: <@%d>",
					*task.TaskName,
					*task.TaskRef,
					*task.AssignedID,
					*task.AssignerID,
				),
			)
		} else if task.AssignedID != nil { // task is assigned, but clear not done by the previous checks

			taskReport.InProgress = append(
				taskReport.InProgress,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Assigned To: <@%d>\n> Assigned By: <@%d>",
					*task.TaskName,
					*task.TaskRef,
					*task.AssignedID,
					*task.AssignerID,
				),
			)
		}
	}
	return taskReport

}
func ParseTaskListWeekly(tasks *[]Task, guildId string) *TaskReport {

	taskReport := &TaskReport{}

	for _, task := range *tasks {

		if *task.Completed { // task is marked compelete by the user and ready for review

			assignedMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignedID))
			assignerMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignerID))

			taskReport.InReview = append(
				taskReport.InReview,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Completed By: %s\n> Waiting for Review: %s",
					*task.TaskName,
					*task.TaskRef,
					assignedMemberName,
					assignerMemberName,
				),
			)
		} else if task.AssignedID != nil { // task is assigned, but clear not done by the previous checks

			assignedMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignedID))
			assignerMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignerID))

			taskReport.InProgress = append(
				taskReport.InProgress,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Assigned To: %s\n> Assigned By: %s",
					*task.TaskName,
					*task.TaskRef,
					assignedMemberName,
					assignerMemberName,
				),
			)
		}
	}
	return taskReport

}

// for time must make some service to convert to the requested timezone, or the timezone the bot is made
func ParseTaskList(tasks *[]Task, guildId string) *TaskReport {

	taskReport := &TaskReport{}

	for _, task := range *tasks {
		if *task.Done { // task is done

			assignedMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignedID))
			assignerMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignerID))

			taskReport.Done = append(
				taskReport.Done,
				fmt.Sprintf("> 📌 **Task Name** %s\n> **Task Ref:** %s\n> Completed By: %s\n> Reviewed By: %s\n> Completed On: %s",
					*task.TaskName,
					*task.TaskRef,
					assignedMemberName,
					assignerMemberName,
					task.FinishedDate.Format("01/02/2006 03:04 PM MST"),
				),
			)
		} else if *task.Completed { // task is marked compelete by the user and ready for review

			assignedMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignedID))
			assignerMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignerID))

			taskReport.InReview = append(
				taskReport.InReview,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Completed By: %s\n> Waiting for Review: %s",
					*task.TaskName,
					*task.TaskRef,
					assignedMemberName,
					assignerMemberName,
				),
			)
		} else if task.AssignedID != nil { // task is assigned, but clear not done by the previous checks

			assignedMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignedID))
			assignerMemberName := GetUserGuildNickname(guildId, strconv.Itoa(*task.AssignerID))

			taskReport.InProgress = append(
				taskReport.InProgress,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Assigned To: %s\n> Assigned By: %s",
					*task.TaskName,
					*task.TaskRef,
					assignedMemberName,
					assignerMemberName,
				),
			)
		} else { // then must be unassigned
			taskReport.Unnassigned = append(
				taskReport.Unnassigned,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n",
					*task.TaskName,
					*task.TaskRef,
				),
			)
		}
	}

	return taskReport
}

func GetUserGuildNickname(guildId string, userId string) string {
	member, err := discord.DiscordSession.GuildMember(guildId, userId)
	if err != nil {
		return ""
	}
	return member.Nick

}

func ValidTaskQuery(s string) bool {
	const base = "milestone"

	if len(s) < len(base) {
		return strings.HasPrefix(base, s)
	}
	if s == base {
		return true
	}

	rest := s[len(base):]

	if ok, _ := regexp.MatchString(`^[0-9]+$`, rest); ok {
		return true
	}

	if ok, _ := regexp.MatchString(`^[0-9]+\/.*`, rest); ok {
		return true
	}
	return false
}

func GetOptionValue(options []*discordgo.ApplicationCommandInteractionDataOption, name string) string {
	for _, opt := range options {
		if opt.Name == name {
			return opt.StringValue()
		}
	}
	return ""
}

func CreateHandleReportAndOutput(success bool, info string, outemdded *discordgo.MessageEmbed, outputId string) *HandleReport {
	return &HandleReport{
		success:    success,
		timestamp:  time.Now(),
		info:       info,
		outputNeed: true,
		outemdded:  outemdded,
		outputId:   outputId,
	}
}

func SetUpProjectInfo(msgInstance *discordgo.InteractionCreate, DB DBClient) (*Project, *HandleReport) {
	if msgInstance.GuildID == "" {
		return nil, CreateHandleReport(false, "Not in a discord server")
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil {
		return nil, CreateHandleReport(false, "channel failed!")
	}

	if channel.ParentID == "" {
		return nil, CreateHandleReport(false, "message must be in a category!")
	}

	project, insertErr := DB.DBGetProject(channel.GuildID, channel.ParentID)

	if insertErr != nil || project == nil {
		return nil, CreateHandleReport(false, "failed to get active project")
	}
	return project, nil
}

func ReportDiscordBotError(err error) {
	discord.DiscordSession.ChannelMessageSend(os.Getenv("OUTPUT_LOG_CHANNEL"), fmt.Sprintf("‼️ <@413398657791164416> Server failure reported!\n Error: %s", err.Error()))
}

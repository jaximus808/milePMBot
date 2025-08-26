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

func CreateAccessRows(msgInstance *discordgo.InteractionCreate, DB DBClient, channel *discordgo.Channel, projectID int) *HandleReport {
	s := discord.DiscordSession

	guildChannels, guildErr := discord.DiscordSession.GuildChannels(msgInstance.GuildID)

	// there should be suffucient checks before this happens but using this to prevent a panic
	if guildErr != nil {
		return CreateHandleReport(false, "No guold Id")
	}

	// gets the first channel that is both text and is within the parent id
	var targetChannel *discordgo.Channel
	for _, ch := range guildChannels {
		if ch.ParentID == channel.ParentID && ch.Type == discordgo.ChannelTypeGuildText {
			targetChannel = ch
			break
		}
	}

	if targetChannel == nil {
		return CreateHandleReport(false, "Not channels could be found in the projects category")
	}

	// we are assuming the discord server will have less than 1000 members, this is a low prior to fix for now, but a future fix if scale is needed

	members, err := s.GuildMembers(msgInstance.GuildID, "", 1000)
	if err != nil {
		return CreateHandleReport(false, err.Error())
	}
	var discordIDs []string
	for _, member := range members {
		// We only care about members who can actually view the channel.
		if member.User.Bot {
			continue
		}
		perms, err := discord.DiscordSession.UserChannelPermissions(member.User.ID, targetChannel.ID)
		if err != nil {
			fmt.Println(err.Error())
		}
		if perms&discordgo.PermissionViewChannel != 0 {
			discordIDs = append(discordIDs, member.User.ID)
		}
	}
	existingUserProfiles, err := DB.DBGetUserProfilesExists(discordIDs)
	if err != nil {
		return CreateHandleReport(false, err.Error())
	}
	// profiles alr exist
	profileMap := make(map[string]int) // discordID -> supabaseUID
	for _, p := range *existingUserProfiles {
		profileMap[strconv.Itoa(*p.DiscordID)] = *p.SupabaseID
	}

	userAccesses := []*UserAccessInsert{}
	newPendingAccess := []*PendingAccessInsert{}

	for _, discordID := range discordIDs {
		supaID, ok := profileMap[discordID]

		if ok {
			userAccesses = append(userAccesses, &UserAccessInsert{
				SupabaseID: &supaID,
				ProjectID:  &projectID,
			})
		} else {
			discordIDNum, err := strconv.Atoi(discordID)
			if err != nil {
				return CreateHandleReport(false, err.Error())
			}
			newPendingAccess = append(newPendingAccess, &PendingAccessInsert{
				DiscordID: &discordIDNum,
				ProjectID: &projectID,
			})
		}
	}
	// now we will create the access and pending access, so when the user logs into supabase for the first time they are given access to the projects tehy dserve access to
	err = DB.DBInsertUserAccess(userAccesses)
	if err != nil {
		return CreateHandleReport(false, err.Error())
	}
	err = DB.DBInsertPendingAccess(newPendingAccess)
	if err != nil {
		return CreateHandleReport(false, err.Error())
	}
	return CreateHandleReport(true, "ok")
}

func ReportDiscordBotError(err error) {
	if err != nil {
		discord.DiscordSession.ChannelMessageSend(os.Getenv("OUTPUT_LOG_CHANNEL"), fmt.Sprintf("‼️ <@413398657791164416> Server failure reported!\n Error: %s", err.Error()))
	} else {
		discord.DiscordSession.ChannelMessageSend(os.Getenv("OUTPUT_LOG_CHANNEL"), fmt.Sprintf("‼️ <@413398657791164416> Server failure reported!\n Error: %s", "Something went wrong!"))
	}
}

func ReportDiscordBotReport(report *HandleReport) {
	discord.DiscordSession.ChannelMessageSend(os.Getenv("OUTPUT_LOG_CHANNEL"), fmt.Sprintf("‼️ <@413398657791164416> Server failure reported!\n Error: %s", report.GetInfo()))
}

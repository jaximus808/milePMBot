package util

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/supabaseutil"
	"github.com/supabase-community/postgrest-go"
)

// TODO: Need to modify the error stuff, maybe also split this up into different files

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
func ptrInt(i int) *int          { return &i }

type HandleReport struct {
	success    bool
	timestamp  time.Time
	info       string
	outputNeed bool
	outemdded  *discordgo.MessageEmbed
	outputId   string
}

type TaskReport struct {
	Unnassigned []string
	InProgress  []string
	InReview    []string
	Done        []string
}

type MilestoneReport struct {
	Upcoming         []string
	CurrentMilestone string
	Previous         []string
}

func (HR *HandleReport) GetInfo() string {
	return HR.info
}
func (HR *HandleReport) GetSuccess() bool {
	return HR.success
}
func (HR *HandleReport) GetTime() time.Time {
	return HR.timestamp
}

func (HR *HandleReport) NeedsOutput() bool {
	return HR.outputNeed
}
func (HR *HandleReport) GetOutputMsg() *discordgo.MessageEmbed {
	return HR.outemdded
}
func (HR *HandleReport) GetOutputId() string {
	return HR.outputId
}
func CreateHandleReport(success bool, info string) *HandleReport {
	return &HandleReport{
		success:    success,
		timestamp:  time.Now(),
		info:       info,
		outputNeed: false,
		outemdded:  nil,
		outputId:   "",
	}
}

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
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Completed By: <%d>\n> Waiting for Review By: <%d>",
					*task.TaskName,
					*task.TaskRef,
					*task.AssignedID,
					*task.AssignerID,
				),
			)
		} else if task.AssignedID != nil { // task is assigned, but clear not done by the previous checks

			taskReport.InProgress = append(
				taskReport.InProgress,
				fmt.Sprintf("> 📌 **Task Name:** %s\n> **Task Ref:** %s\n> Assigned To: <%d>\n> Assigned By: <%d>",
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

func SetUpProjectInfo(msgInstance *discordgo.InteractionCreate) (*Project, *HandleReport) {
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

	project, insertErr := DBGetProject(channel.GuildID, channel.ParentID)

	if insertErr != nil || project == nil {
		return nil, CreateHandleReport(false, "failed to get active project")
	}
	return project, nil
}

/**
Active Project Methods
*/

func DBCreateActiveProject(guildId int, pchannelId int, projectId int) (*ActiveProject, error) {
	var selectedActiveProject ActiveProject
	insertedActiveProject := ActiveProjectInsert{
		GuildID:    &guildId,
		PChannelID: &pchannelId,
		ProjectID:  &projectId,
	}
	res, _, err := supabaseutil.Client.From("ActiveProjects").Insert(insertedActiveProject, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedActiveProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &selectedActiveProject, nil
}

func DBGetActiveProject(guildId string, pchannelId string) (*ActiveProject, error) {
	var selectedActiveProject ActiveProject

	res, _, err := supabaseutil.Client.From("ActiveProjects").Select("*", "", false).Eq("pChannelId", pchannelId).Eq("guildId", guildId).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedActiveProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &selectedActiveProject, nil
}

/**
Projects Methods
*/

func DBCreateProject(guildId int, pchannelId int, outoutChannelId int, desc string) (*Project, error) {
	var selectedProject Project
	insertedProject := ProjectInsert{
		GuildID:       &guildId,
		PChannelID:    &pchannelId,
		SprintMsg:     ptrString("Sprint Message: "),
		OutputChannel: &outoutChannelId,
		Desc:          &desc,
	}
	res, _, err := supabaseutil.Client.From("Projects").Insert(insertedProject, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &selectedProject, nil
}

func DBGetProjectId(projectId int) (*Project, error) {
	var selectedProject Project

	//gets the active project
	projectRes, _, err := supabaseutil.Client.From("Projects").Select("*", "", false).Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(projectRes, &selectedProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedProject, nil
}

func DBGetProject(guildId string, pchannelId string) (*Project, error) {
	var selectedProject Project
	var activeProject ActiveProject

	//gets the active project
	activeRes, _, err := supabaseutil.Client.From("ActiveProjects").Select("*", "", false).Eq("pChannelId", pchannelId).Eq("guildId", guildId).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(activeRes, &activeProject)
	if err != nil {
		log.Printf("Active Error unmarshaling response: %v", err)
		return nil, err
	}
	//gets the active project
	projectRes, _, err := supabaseutil.Client.From("Projects").Select("*", "", false).Eq("id", strconv.Itoa(*activeProject.ProjectID)).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(projectRes, &selectedProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedProject, nil
}

func DBUpdateProjectOutputChannel(projectId int, outputId int) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		OutputChannel: &outputId,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateProjectDescription(projectId int, desc string) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		Desc: &desc,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateProjectSprintDesc(projectId int, desc string) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		SprintMsg: &desc,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateProjectSprints(projectId int, enabled bool) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		SprintEnabled: &enabled,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateProjectSprintDuration(projectId int, sprintDuration int) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		SprintN: &sprintDuration,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateResetSprintDuration(projectId int) (*Project, error) {
	var newProject Project
	// push back a day for cushion
	updatedProject := ProjectUpdate{
		LastPingAt: time.Now().Add(-24 * time.Hour),
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBUpdateProjectPings(projectId int, enabled bool) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		SprintPing: &enabled,
		LastPingAt: time.Now(),
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(updatedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newProject, nil
}

func DBGetAllPingProjects() (*([]Project), error) {
	var projects []Project
	res, _, err := supabaseutil.Client.From("Projects").Select("*,ActiveProjects!inner(project_id)", "", false).Eq("sprint_enabled", "TRUE").Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &projects)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &projects, nil
}

// Milestone commands

func DBCreateMilestone(projectId int, msName string, msDeadline *time.Time, msDes string) (*Milestone, error) {
	var selectedMilestone Milestone
	insertedMilestone := MilestoneInsert{
		ProjectID:   &projectId,
		DisplayName: &msName,
		DueDate:     msDeadline,
		Description: &msDes,
	}
	res, _, err := supabaseutil.Client.From("Milestones").Insert(insertedMilestone, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &selectedMilestone, nil
}

func DBGetMilestoneWithId(milestoneID int) (*Milestone, error) {
	var selectedMilestone Milestone
	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("id", strconv.Itoa(milestoneID)).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &selectedMilestone, nil
}

func DBGetNextMilestone(projectId int, curMilstone *Milestone) (*Milestone, error) {
	var nextMilestone Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = true

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Gt("due_date", curMilstone.DueDate.Format(time.RFC3339)).Order("due_date", orderOptions).Limit(1, "").Single().Execute()

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &nextMilestone)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &nextMilestone, nil
}

func DBGetPrevMilestone(projectId int, curMilstone *Milestone) (*Milestone, error) {
	var nextMilestone Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = false

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Lt("due_date", curMilstone.DueDate.Format(time.RFC3339)).Order("due_date", orderOptions).Limit(1, "").Single().Execute()

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &nextMilestone)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &nextMilestone, nil
}
func DBUpdateCurrentMilestone(projectId int, milestoneId int) (*Project, error) {
	var updatedProject Project
	insertedProject := ProjectInsert{
		CurrentMID: &milestoneId,
	}
	res, _, err := supabaseutil.Client.From("Projects").Update(insertedProject, "representation", "").Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("BROOO Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &updatedProject)
	if err != nil {
		log.Printf(" WTFFF Error unmarshaling response: %v", err)
		return nil, err
	}
	return &updatedProject, nil
}

func DBMilestoneExistDate(projectId int, deadline *time.Time) bool {
	_, _, err := supabaseutil.Client.From("Milestones").Select("id", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("due_date", deadline.Format(time.RFC3339)).Single().Execute()

	return err == nil
}

func DBGetMilestoneListDescending(projectId int) (*[]Milestone, error) {
	var milestoneList []Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = false

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Order("due_date", orderOptions).Execute()

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &milestoneList)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &milestoneList, nil
}

/*
*

	Role comamands
*/
func DBCreateRole(projectId int, userId int, roleLevel int) (*Role, error) {
	var selectedRole Role
	insertedRole := RoleInsert{
		ProjectID: &projectId,
		DiscordID: &userId,
		RoleLevel: roleLevel,
	}
	res, _, err := supabaseutil.Client.From("Roles").Insert(insertedRole, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedRole)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedRole, nil
}

func DBGetRole(projectId int, userId string) (*Role, error) {
	var selectedRole Role

	res, _, err := supabaseutil.Client.From("Roles").Select("*", "", false).Eq("discord_id", userId).Eq("project_id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedRole)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedRole, nil
}
func DBDeleteRole(roleId int) error {
	_, _, err := supabaseutil.Client.From("Roles").Delete("*", "").Eq("id", strconv.Itoa(roleId)).Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return err
	}

	return nil
}

/*
	Task commands
*/

func DBCreateTasks(projectId int, task_name string, desc string, milestone_id int) (*Task, error) {
	var newTask Task
	insertedTask := TaskInsert{
		TaskName:    &task_name,
		TaskRef:     ptrString(fmt.Sprintf("milestone%d/%s", milestone_id, strings.ReplaceAll(task_name, " ", "_"))),
		ProjectID:   &projectId,
		Desc:        &desc,
		MilestoneID: &milestone_id,
	}
	res, _, err := supabaseutil.Client.From("Tasks").Insert(insertedTask, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &newTask, nil
}
func DBGetTask(projectId int, taskRef string) (*Task, error) {
	var selectedTask Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedTask, nil
}

func DBGetTasksForMilestone(projectId int, milestoneId int) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedTasks, nil
}
func DBGetTasksForMilestoneAndAssignedUser(projectId int, milestoneId int, userId string) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Eq("assigned_id", userId).Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedTasks, nil
}

func DBAssignTasksDueDate(projectId int, taskRef string, assignerId int, assignedId int, dueDate *time.Time) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		AssignedID: &assignedId,
		DueDate:    dueDate,
		AssignerID: &assignedId,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Single().Execute()
	if err != nil {
		log.Printf("MEOWW ?Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newTask, nil
}

func DBAssignTasksStoryPoints(projectId int, taskRef string, assignerId int, assignedId int, storyPoint int) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		AssignedID:  &assignedId,
		StoryPoints: &storyPoint,
		AssignerID:  &assignerId,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Eq("done", "false").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newTask, nil
}
func DBUpdateTaskRecentProgress(taskId int, completed bool) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		Completed: &completed,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("id", strconv.Itoa(taskId)).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newTask, nil
}

func DBGetInProgressAndCompetedTask(projectId int, milestoneId int) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Not("assigned_id", "is", "NULL").Eq("done", "FALSE").Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &selectedTasks, nil
}

// used for approve, lowky this shit is so fucking bad
func DBTaskMarkComplete(projectId int, taskRef string, finishedDate *time.Time) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		Done:         ptrBool(true),
		FinishedDate: finishedDate,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	return &newTask, nil
}

func DBGetSimillarTasksAssigned(discordId string, taskRefQuery string, isAssigner bool, projectId int) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error
	if isAssigner {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigner_id", discordId).Ilike("task_ref", taskRefQuery+"%").Execute()
	} else {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigned_id", discordId).Ilike("task_ref", taskRefQuery+"%").Execute()
	}

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &tasksMatch, nil

}

func DBGetUnassignedTasks(discordId string, taskRefQuery string, projectId int) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error

	res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Is("assigned_id", "NULL").Ilike("task_ref", taskRefQuery+"%").Execute()

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &tasksMatch, nil
}

func DBGetTasksAndSpecifyDC(discordId string, taskRefQuery string, isAssigner bool, projectId int, done bool, complete bool) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error
	if isAssigner {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigner_id", discordId).Eq("completed", strings.ToUpper(strconv.FormatBool(done))).Eq("done", strings.ToUpper(strconv.FormatBool(complete))).Ilike("task_ref", taskRefQuery+"%").Execute()
	} else {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigned_id", discordId).Eq("completed", strings.ToUpper(strconv.FormatBool(done))).Eq("done", strings.ToUpper(strconv.FormatBool(complete))).Ilike("task_ref", taskRefQuery+"%").Execute()
	}

	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &tasksMatch, nil
}

//Progress commands

func DBCreateProgress(task_id int, desc string, completed bool) (*Progress, error) {
	var newProgress Progress
	insertedProgress := ProgressInsert{
		TaskID:    &task_id,
		Desc:      &desc,
		Completed: &completed,
	}
	res, _, err := supabaseutil.Client.From("Progress").Insert(insertedProgress, false, "", "representation", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &newProgress)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	return &newProgress, nil
}

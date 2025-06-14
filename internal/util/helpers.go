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
	outmsg     string
	outputId   string
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
func (HR *HandleReport) GetOutputMsg() string {
	return HR.outmsg
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
		outmsg:     "",
		outputId:   "",
	}
}
func ValidTaskQuery(s string) bool {
	const base = "milestone"

	// 1) still typing the word “milestone”?
	if len(s) < len(base) {
		return strings.HasPrefix(base, s)
	}
	// 2) they've typed all of “milestone” and maybe more
	if s == base {
		return true
	}

	rest := s[len(base):] // what comes after "milestone"

	// 3) digits only → still a valid milestone ID
	if ok, _ := regexp.MatchString(`^[0-9]+$`, rest); ok {
		return true
	}
	// 4) digits + "/" + anything → now we're querying sub-items under that milestone
	if ok, _ := regexp.MatchString(`^[0-9]+\/.*`, rest); ok {
		return true
	}
	// otherwise bail out
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

func CreateHandleReportAndOutput(success bool, info string, outMsg string, outputId string) *HandleReport {
	return &HandleReport{
		success:    success,
		timestamp:  time.Now(),
		info:       info,
		outputNeed: true,
		outmsg:     outMsg,
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

func DBCreateActiveProject(guildId string, pchannelId string, projectId int) (*ActiveProject, error) {
	var selectedActiveProject ActiveProject
	insertedActiveProject := ActiveProjectInsert{
		GuildID:    ptrString(guildId),
		PChannelID: ptrString(pchannelId),
		ProjectID:  ptrInt(projectId),
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

func DBCreateProject(guildId string, pchannelId string, outoutChannelId string, desc string) (*Project, error) {
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

func DBMilestoneExistDate(projectId int, deadline *time.Time) bool {
	_, _, err := supabaseutil.Client.From("Milestones").Select("id", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("due_date", deadline.Format(time.RFC3339)).Single().Execute()

	return err == nil
}

/*
*

	Role comamands
*/
func DBCreateRole(projectId int, userId string, roleLevel int) (*Role, error) {
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
func DBAssignTasksDueDate(projectId int, taskRef string, assignerId string, assignedId string, dueDate *time.Time) (*Task, error) {
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

func DBAssignTasksStoryPoints(projectId int, taskRef string, assignerId string, assignedId string, storyPoint int) (*Task, error) {
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

package util

import (
	"encoding/json"
	"log"
	"strconv"
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
	success   bool
	timestamp time.Time
	info      string
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

func CreateHandleReport(success bool, info string) *HandleReport {
	return &HandleReport{
		success:   success,
		timestamp: time.Now(),
		info:      info,
	}
}

func SetUpProjectInfo(msgInstance *discordgo.MessageCreate) (*Project, *HandleReport) {
	if msgInstance.GuildID == "" {
		return nil, CreateHandleReport(false, "Not in a discord server")
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil {
		return nil, CreateHandleReport(false, "channel failed!")
	}

	if channel.Type != 4 {
		return nil, CreateHandleReport(false, "message must be in a category!")
	}

	project, insertErr := DBGetProject(channel.GuildID, channel.ParentID)

	if insertErr != nil || project == nil {
		return nil, CreateHandleReport(false, "failed to make project")
	}
	return project, nil
}

/**
Projects Methods
*/

func DBCreateProject(guildId string, pchannelId string, outoutChannelId string) (*Project, error) {
	var selectedProject Project
	insertedProject := ProjectInsert{
		GuildID:       ptrString(guildId),
		PChannelID:    ptrString(pchannelId),
		SprintMsg:     ptrString("Sprint Message: "),
		OutputChannel: ptrString(outoutChannelId),
	}
	res, _, err := supabaseutil.Client.From("Projects").Insert(insertedProject, false, "", "Projects", "*").Single().Execute()
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

func DBGetProject(guildId string, pchannelId string) (*Project, error) {
	var selectedProject Project
	var activeProject ActiveProject

	//gets the active project
	activeRes, _, err := supabaseutil.Client.From("ActiveProjects").Select("*", "exact", false).Eq("pChannelId", pchannelId).Eq("guildId", guildId).Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(activeRes, &activeProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}

	//gets the active project
	projectRes, _, err := supabaseutil.Client.From("Projects").Select("*", "exact", false).Eq("id", strconv.Itoa(activeProject.ID)).Execute()
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
	res, _, err := supabaseutil.Client.From("Projects").Update(insertedProject, "Projects", "").Single().Execute()
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return nil, err
	}
	err = json.Unmarshal(res, &updatedProject)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
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
	res, _, err := supabaseutil.Client.From("Milestones").Insert(insertedMilestone, false, "", "Milestones", "*").Single().Execute()
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
	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "exact", false).Eq("id", strconv.Itoa(milestoneID)).Single().Execute()
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

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "exact", false).Eq("project_id", strconv.Itoa(projectId)).Gt("due_date", strconv.Itoa(int(curMilstone.DueDate.Unix()))).Order("due_date", orderOptions).Limit(1, "Milestones").Single().Execute()

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &nextMilestone)
	if err != nil {
		return nil, err
	}
	return &nextMilestone, nil
}

func DBGetPrevMilestone(projectId int, curMilstone *Milestone) (*Milestone, error) {
	var prevMilestone Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = false

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "exact", false).Eq("project_id", strconv.Itoa(projectId)).Lt("due_date", strconv.Itoa(int(curMilstone.DueDate.Unix()))).Order("due_date", orderOptions).Limit(1, "Milestones").Single().Execute()

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &prevMilestone)
	if err != nil {
		return nil, err
	}
	return &prevMilestone, nil
}
func DBCreateRole(projectId int, userId string, roleLevel int) (*Role, error) {
	var selectedRole Role
	insertedRole := RoleInsert{
		ProjectID: &projectId,
		DiscordID: &userId,
		RoleLevel: roleLevel,
	}
	res, _, err := supabaseutil.Client.From("Roles").Insert(insertedRole, false, "", "Roles", "*").Single().Execute()
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

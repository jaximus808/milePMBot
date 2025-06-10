package util

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

func ptrString(s string) *string { return &s }
func ptrBool(b bool) *bool       { return &b }
func ptrInt(i int) *int          { return &i }

type HandleReport struct {
	success   bool
	timestamp time.Time
	info      string
}

func CreateHandleReport(success bool, info string) *HandleReport {
	return &HandleReport{
		success:   success,
		timestamp: time.Now(),
		info:      info,
	}
}

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

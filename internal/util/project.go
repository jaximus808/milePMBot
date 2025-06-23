package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

/**
Projects Methods
*/

func (s SupaDB) DBCreateProject(guildId int, pchannelId int, outoutChannelId int, desc string) (*Project, error) {
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

		return nil, err
	}
	err = json.Unmarshal(res, &selectedProject)
	if err != nil {

		return nil, err
	}
	return &selectedProject, nil
}

func (s SupaDB) DBGetProjectId(projectId int) (*Project, error) {
	var selectedProject Project

	//gets the active project
	projectRes, _, err := supabaseutil.Client.From("Projects").Select("*", "", false).Eq("id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(projectRes, &selectedProject)
	if err != nil {

		return nil, err
	}

	return &selectedProject, nil
}
func (s SupaDB) DBGetProjectRef(projectRef string) (*Project, error) {
	var selectedProject Project

	//gets the active project
	projectRes, _, err := supabaseutil.Client.From("Projects").Select("*", "", false).Eq("project_ref", projectRef).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(projectRes, &selectedProject)
	if err != nil {

		return nil, err
	}

	return &selectedProject, nil
}

func (s SupaDB) DBGetProject(guildId string, pchannelId string) (*Project, error) {
	var selectedProject Project
	var activeProject ActiveProject

	//gets the active project
	activeRes, _, err := supabaseutil.Client.From("ActiveProjects").Select("*", "", false).Eq("pChannelId", pchannelId).Eq("guildId", guildId).Single().Execute()
	if err != nil {

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

		return nil, err
	}
	err = json.Unmarshal(projectRes, &selectedProject)
	if err != nil {

		return nil, err
	}

	return &selectedProject, nil
}
func (s SupaDB) DBUpdateProjectRef(projectId int) (*Project, error) {
	var newProject Project
	updatedProject := ProjectUpdate{
		ProjectRef: ptrString(fmt.Sprintf("project%d", projectId)),
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

func (s SupaDB) DBUpdateProjectOutputChannel(projectId int, outputId int) (*Project, error) {
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

func (s SupaDB) DBUpdateProjectDescription(projectId int, desc string) (*Project, error) {
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

func (s SupaDB) DBUpdateProjectSprintDesc(projectId int, desc string) (*Project, error) {
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

func (s SupaDB) DBUpdateProjectSprints(projectId int, enabled bool) (*Project, error) {
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

func (s SupaDB) DBUpdateProjectSprintDuration(projectId int, sprintDuration int) (*Project, error) {
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

func (s SupaDB) DBUpdateResetSprintDuration(projectId int) (*Project, error) {
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

func (s SupaDB) DBUpdateProjectPings(projectId int, enabled bool) (*Project, error) {
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

func (s SupaDB) DBGetAllPingProjects() (*([]Project), error) {
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

/*
*

	Role comamands
*/
func (s SupaDB) DBCreateRole(projectId int, userId int, roleLevel int) (*Role, error) {
	var selectedRole Role
	insertedRole := RoleInsert{
		ProjectID: &projectId,
		DiscordID: &userId,
		RoleLevel: roleLevel,
	}
	res, _, err := supabaseutil.Client.From("Roles").Insert(insertedRole, false, "", "representation", "").Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedRole)
	if err != nil {

		return nil, err
	}

	return &selectedRole, nil
}

func (s SupaDB) DBGetRole(projectId int, userId string) (*Role, error) {
	var selectedRole Role

	res, _, err := supabaseutil.Client.From("Roles").Select("*", "", false).Eq("discord_id", userId).Eq("project_id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedRole)
	if err != nil {

		return nil, err
	}

	return &selectedRole, nil
}
func (s SupaDB) DBDeleteRole(roleId int) error {
	_, _, err := supabaseutil.Client.From("Roles").Delete("*", "").Eq("id", strconv.Itoa(roleId)).Execute()
	if err != nil {

		return err
	}

	return nil
}

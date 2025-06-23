package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
	"github.com/supabase-community/postgrest-go"
)

// Milestone commands

func (s SupaDB) DBCreateMilestone(projectId int, msName string, msDeadline *time.Time, msDes string) (*Milestone, error) {
	var selectedMilestone Milestone
	insertedMilestone := MilestoneInsert{
		ProjectID:   &projectId,
		DisplayName: &msName,
		DueDate:     msDeadline,
		Description: &msDes,
	}
	res, _, err := supabaseutil.Client.From("Milestones").Insert(insertedMilestone, false, "", "representation", "").Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {

		return nil, err
	}
	return &selectedMilestone, nil
}

func (s SupaDB) DBUpdateMilestoneRef(milestoneId int, projectRef string) (*Milestone, error) {
	var selectedMilestone Milestone
	updatedMilestone := MilestoneInsert{
		MilestoneRef: ptrString(fmt.Sprintf("%s/milestone%d", projectRef, milestoneId)),
	}
	res, _, err := supabaseutil.Client.From("Milestones").Update(updatedMilestone, "representation", "").Eq("id", strconv.Itoa(milestoneId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {

		return nil, err
	}
	return &selectedMilestone, nil
}

func (s SupaDB) DBGetMilestoneWithId(milestoneID int) (*Milestone, error) {
	var selectedMilestone Milestone
	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("id", strconv.Itoa(milestoneID)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {

		return nil, err
	}
	return &selectedMilestone, nil
}

func (s SupaDB) DBGetMilestoneWithRef(milestoneRef string, projectId int) (*Milestone, error) {
	var selectedMilestone Milestone
	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("milestone_ref", milestoneRef).Eq("project_id", strconv.Itoa(projectId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedMilestone)
	if err != nil {

		return nil, err
	}
	return &selectedMilestone, nil
}
func (s SupaDB) DBDeleteMilestone(milestoneId int) error {
	_, _, err := supabaseutil.Client.From("Milestones").Delete("*", "").Eq("id", strconv.Itoa(milestoneId)).Execute()
	return err
}
func (s SupaDB) DBGetNextMilestone(projectId int, curMilstone *Milestone) (*Milestone, error) {
	var nextMilestone Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = true

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Gt("due_date", curMilstone.DueDate.Format(time.RFC3339)).Order("due_date", orderOptions).Limit(1, "").Single().Execute()

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &nextMilestone)
	if err != nil {

		return nil, err
	}
	return &nextMilestone, nil
}

func (s SupaDB) DBGetPrevMilestone(projectId int, curMilstone *Milestone) (*Milestone, error) {
	var nextMilestone Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = false

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Lt("due_date", curMilstone.DueDate.Format(time.RFC3339)).Order("due_date", orderOptions).Limit(1, "").Single().Execute()

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &nextMilestone)
	if err != nil {

		return nil, err
	}
	return &nextMilestone, nil
}
func (s SupaDB) DBUpdateCurrentMilestone(projectId int, milestoneId int) (*Project, error) {
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

func (s SupaDB) DBMilestoneExistDate(projectId int, deadline *time.Time) bool {
	_, _, err := supabaseutil.Client.From("Milestones").Select("id", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("due_date", deadline.Format(time.RFC3339)).Single().Execute()

	return err == nil
}

func (s SupaDB) DBGetMilestoneListDescending(projectId int) (*[]Milestone, error) {
	var milestoneList []Milestone
	orderOptions := &postgrest.DefaultOrderOpts

	orderOptions.Ascending = false

	res, _, err := supabaseutil.Client.From("Milestones").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Order("due_date", orderOptions).Execute()

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &milestoneList)
	if err != nil {

		return nil, err
	}
	return &milestoneList, nil
}

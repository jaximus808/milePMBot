package util

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jaximus808/milePMBot/internal/supabaseutil"
)

/*
	Task commands
*/

func (s SupaDB) DBCreateTasks(projectId int, task_name string, desc string, milestone_id int) (*Task, error) {
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

		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {

		return nil, err
	}

	return &newTask, nil
}
func (s SupaDB) DBGetTask(projectId int, taskRef string) (*Task, error) {
	var selectedTask Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedTask)
	if err != nil {

		return nil, err
	}

	return &selectedTask, nil
}

func (s SupaDB) DBGetTasksForMilestone(projectId int, milestoneId int) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {

		return nil, err
	}

	return &selectedTasks, nil
}
func (s SupaDB) DBGetTasksForMilestoneAndAssignedUser(projectId int, milestoneId int, userId string) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Eq("assigned_id", userId).Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {

		return nil, err
	}

	return &selectedTasks, nil
}

func (s SupaDB) DBAssignTasksDueDate(projectId int, taskRef string, assignerId int, assignedId int, dueDate *time.Time) (*Task, error) {
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

		return nil, err
	}
	return &newTask, nil
}

func (s SupaDB) DBAssignTasksStoryPoints(projectId int, taskRef string, assignerId int, assignedId int, storyPoint int) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		AssignedID:  &assignedId,
		StoryPoints: &storyPoint,
		AssignerID:  &assignerId,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Eq("done", "false").Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {

		return nil, err
	}
	return &newTask, nil
}
func (s SupaDB) DBUpdateTaskRecentProgress(taskId int, completed bool) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		Completed: &completed,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("id", strconv.Itoa(taskId)).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {

		return nil, err
	}
	return &newTask, nil
}

func (s SupaDB) DBGetInProgressAndCompetedTask(projectId int, milestoneId int) (*[]Task, error) {
	var selectedTasks []Task
	res, _, err := supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("milestone_id", strconv.Itoa(milestoneId)).Not("assigned_id", "is", "NULL").Eq("done", "FALSE").Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &selectedTasks)
	if err != nil {

		return nil, err
	}

	return &selectedTasks, nil
}

// used for approve, lowky this shit is so fucking bad
func (s SupaDB) DBTaskMarkComplete(projectId int, taskRef string, finishedDate *time.Time) (*Task, error) {
	var newTask Task
	updatedTask := TaskUpdate{
		Done:         ptrBool(true),
		FinishedDate: finishedDate,
	}

	res, _, err := supabaseutil.Client.From("Tasks").Update(updatedTask, "representation", "").Eq("project_id", strconv.Itoa(projectId)).Eq("task_ref", taskRef).Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &newTask)
	if err != nil {

		return nil, err
	}
	return &newTask, nil
}

func (s SupaDB) DBGetSimillarTasksAssigned(discordId string, taskRefQuery string, isAssigner bool, projectId int) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error
	if isAssigner {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigner_id", discordId).Ilike("task_ref", taskRefQuery+"%").Execute()
	} else {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigned_id", discordId).Ilike("task_ref", taskRefQuery+"%").Execute()
	}

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {

		return nil, err
	}

	return &tasksMatch, nil

}

func (s SupaDB) DBGetUnassignedTasks(discordId string, taskRefQuery string, projectId int) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error

	res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Is("assigned_id", "NULL").Ilike("task_ref", taskRefQuery+"%").Execute()

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {

		return nil, err
	}

	return &tasksMatch, nil
}

func (s SupaDB) DBGetTasksAndSpecifyDC(discordId string, taskRefQuery string, isAssigner bool, projectId int, done bool, complete bool) (*[]Task, error) {
	var tasksMatch []Task
	var res []byte
	var err error
	if isAssigner {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigner_id", discordId).Eq("completed", strings.ToUpper(strconv.FormatBool(done))).Eq("done", strings.ToUpper(strconv.FormatBool(complete))).Ilike("task_ref", taskRefQuery+"%").Execute()
	} else {
		res, _, err = supabaseutil.Client.From("Tasks").Select("*", "", false).Eq("project_id", strconv.Itoa(projectId)).Eq("assigned_id", discordId).Eq("completed", strings.ToUpper(strconv.FormatBool(done))).Eq("done", strings.ToUpper(strconv.FormatBool(complete))).Ilike("task_ref", taskRefQuery+"%").Execute()
	}

	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &tasksMatch)
	if err != nil {

		return nil, err
	}

	return &tasksMatch, nil
}

//Progress commands

func (s SupaDB) DBCreateProgress(task_id int, desc string, completed bool) (*Progress, error) {
	var newProgress Progress
	insertedProgress := ProgressInsert{
		TaskID:    &task_id,
		Desc:      &desc,
		Completed: &completed,
	}
	res, _, err := supabaseutil.Client.From("Progress").Insert(insertedProgress, false, "", "representation", "").Single().Execute()
	if err != nil {

		return nil, err
	}
	err = json.Unmarshal(res, &newProgress)
	if err != nil {

		return nil, err
	}

	return &newProgress, nil
}

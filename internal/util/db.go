package util

import "time"

type DBClient interface {

	// active project table function

	DBCreateActiveProject(guildId int, pchannelId int, projectId int) (*ActiveProject, error)

	DBGetActiveProject(guildId string, pchannelId string) (*ActiveProject, error)

	DBUpdateProjectId(projectId int, guildId int, pchannelId int) (*ActiveProject, error)

	DBEndActiveProject(projectId int) error

	// projects

	DBCreateProject(guildId int, pchannelId int, outoutChannelId int, desc string) (*Project, error)

	DBGetProjectId(projectId int) (*Project, error)

	DBGetProjectRef(projectRef string) (*Project, error)

	DBGetProject(guildId string, pchannelId string) (*Project, error)

	DBUpdateProjectRef(projectId int) (*Project, error)

	DBUpdateProjectOutputChannel(projectId int, outputId int) (*Project, error)

	DBUpdateProjectDescription(projectId int, desc string) (*Project, error)

	DBUpdateProjectSprintDesc(projectId int, desc string) (*Project, error)

	DBUpdateProjectSprints(projectId int, enabled bool) (*Project, error)

	DBUpdateProjectSprintDuration(projectId int, sprintDuration int) (*Project, error)

	DBUpdateResetSprintDuration(projectId int) (*Project, error)

	DBUpdateProjectPings(projectId int, enabled bool) (*Project, error)

	DBGetAllProjects() (*([]Project), error)

	DBGetAllPingProjects() (*([]Project), error)

	DBCreateRole(projectId int, userId int, roleLevel int) (*Role, error)

	DBGetRole(projectId int, userId string) (*Role, error)

	DBDeleteRole(roleId int) error

	DBDeleteProject(projectId int) error

	// milestones

	DBCreateMilestone(projectId int, msName string, msDeadline *time.Time, msDes string) (*Milestone, error)

	DBUpdateMilestoneRef(milestoneId int, projectRef string) (*Milestone, error)

	DBGetMilestoneWithId(milestoneID int) (*Milestone, error)

	DBGetMilestoneWithRef(milestoneRef string, projectId int) (*Milestone, error)

	DBDeleteMilestone(milestoneId int) error

	DBGetNextMilestone(projectId int, curMilstone *Milestone) (*Milestone, error)

	DBGetPrevMilestone(projectId int, curMilstone *Milestone) (*Milestone, error)

	DBUpdateCurrentMilestone(projectId int, milestoneId int) (*Project, error)

	DBMilestoneExistDate(projectId int, deadline *time.Time) bool

	DBGetMilestoneListDescending(projectId int) (*[]Milestone, error)

	// task
	DBCreateTasks(projectId int, task_name string, desc string, milestone_id int) (*Task, error)

	DBGetTask(projectId int, taskRef string) (*Task, error)

	DBGetTasksForMilestone(projectId int, milestoneId int) (*[]Task, error)

	DBGetTasksForMilestoneAndAssignedUser(projectId int, milestoneId int, userId string) (*[]Task, error)

	DBAssignTasksDueDate(projectId int, taskRef string, assignerId int, assignedId int, dueDate *time.Time) (*Task, error)

	DBAssignTasksStoryPoints(projectId int, taskRef string, assignerId int, assignedId int, storyPoint int) (*Task, error)

	DBUpdateTaskRecentProgress(taskId int, completed bool) (*Task, error)

	DBGetInProgressAndCompetedTask(projectId int, milestoneId int) (*[]Task, error)

	DBTaskMarkComplete(projectId int, taskRef string, finishedDate *time.Time) (*Task, error)

	DBGetSimillarTasksAssigned(discordId string, taskRefQuery string, isAssigner bool, projectId int) (*[]Task, error)

	DBGetUnassignedTasks(discordId string, taskRefQuery string, projectId int) (*[]Task, error)

	DBGetTasksAndSpecifyDC(discordId string, taskRefQuery string, isAssigner bool, projectId int, done bool, complete bool) (*[]Task, error)

	DBCreateProgress(task_id int, desc string, completed bool) (*Progress, error)
}

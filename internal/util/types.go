package util

import "time"

type RoleLevel int

const (
	DefaultRole RoleLevel = iota
	LeadRole
	AdminRole
)

// ActiveProjects
type ActiveProject struct {
	CreatedAt  time.Time `json:"created_at"`
	GuildID    *string   `json:"guildId"`
	ID         int       `json:"id"`
	PChannelID *string   `json:"pChannelId"`
	ProjectID  *int      `json:"project_id"`
}

type ActiveProjectInsert struct {
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	GuildID    *string    `json:"guildId,omitempty"`
	ID         *int       `json:"id,omitempty"`
	PChannelID *string    `json:"pChannelId,omitempty"`
	ProjectID  *int       `json:"project_id,omitempty"`
}

type ActiveProjectUpdate = ActiveProjectInsert

// Milestones
type Milestone struct {
	CreatedAt    time.Time  `json:"created_at"`
	Description  *string    `json:"description"`
	DisplayName  *string    `json:"display_name"`
	DueDate      *time.Time `json:"due_date"`
	ID           int        `json:"id"`
	MilestoneRef *string    `json:"milestone_ref"`
	ProjectID    *int       `json:"project_id"`
}

type MilestoneInsert struct {
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Description  *string    `json:"description,omitempty"`
	DisplayName  *string    `json:"display_name,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	ID           *int       `json:"id,omitempty"`
	MilestoneRef *string    `json:"milestone_ref,omitempty"`
	ProjectID    *int       `json:"project_id,omitempty"`
}

type MilestoneUpdate = MilestoneInsert

// Progress
type Progress struct {
	CreatedAt time.Time `json:"created_at"`
	Desc      *string   `json:"desc"`
	ID        int       `json:"id"`
	TaskID    *int      `json:"task_id"`
	Completed *bool     `json:"completed,omitempty"`
}

type ProgressInsert struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Desc      *string    `json:"desc,omitempty"`
	ID        *int       `json:"id,omitempty"`
	TaskID    *int       `json:"task_id,omitempty"`
	Completed *bool      `json:"completed,omitempty"`
}

type ProgressUpdate = ProgressInsert

// Projects
type Project struct {
	Completed     *bool     `json:"completed"`
	CreatedAt     time.Time `json:"created_at"`
	CurrentMID    *int      `json:"current_mid"`
	Desc          *string   `json:"desc"`
	GuildID       *string   `json:"guild_id"`
	ID            int       `json:"id"`
	OutputChannel *string   `json:"output_channel"`
	PChannelID    *string   `json:"pchannel_id"`
	SprintEnabled *bool     `json:"sprint_enabled"`
	SprintInt     *int      `json:"sprint_int"`
	SprintMsg     *string   `json:"sprint_msg"`
	SprintN       *int      `json:"sprint_n"`
	SprintPing    *bool     `json:"sprint_ping"`
}

type ProjectInsert struct {
	Completed     *bool      `json:"completed,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	CurrentMID    *int       `json:"current_mid,omitempty"`
	Desc          *string    `json:"desc,omitempty"`
	GuildID       *string    `json:"guild_id,omitempty"`
	ID            *int       `json:"id,omitempty"`
	OutputChannel *string    `json:"output_channel,omitempty"`
	PChannelID    *string    `json:"pchannel_id,omitempty"`
	SprintEnabled *bool      `json:"sprint_enabled,omitempty"`
	SprintInt     *int       `json:"sprint_int,omitempty"`
	SprintMsg     *string    `json:"sprint_msg,omitempty"`
	SprintN       *int       `json:"sprint_n,omitempty"`
	SprintPing    *bool      `json:"sprint_ping,omitempty"`
}

type ProjectUpdate = ProjectInsert

// Roles
type Role struct {
	CreatedAt time.Time `json:"created_at"`
	DiscordID *string   `json:"discord_id"`
	ID        int       `json:"id"`
	ProjectID *int      `json:"project_id"`
	RoleLevel int       `json:"role_level"`
}

type RoleInsert struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	DiscordID *string    `json:"discord_id,omitempty"`
	ID        *int       `json:"id,omitempty"`
	ProjectID *int       `json:"project_id,omitempty"`
	RoleLevel int        `json:"role_level"`
}

type RoleUpdate = RoleInsert

// Tasks
type Task struct {
	AssignedID       *string    `json:"assigned_id"`
	AssignerID       *string    `json:"assigner_id"`
	CreatedAt        time.Time  `json:"created_at"`
	Desc             *string    `json:"desc"`
	Done             *bool      `json:"done"`
	DueDate          *time.Time `json:"due_date"`
	ID               int        `json:"id"`
	MilestoneID      *int       `json:"milestone_id"`
	TaskName         *string    `json:"task_name"`
	TaskRef          *string    `json:"task_ref"`
	ProjectID        *int       `json:"project_id,omitempty"`
	StoryPoints      *int       `json:"story_points,omitempty"`
	RecentProgressId *int       `json:"recent_progress"`
	FinishedDate     *time.Time `json:"finished_date,omitempty"`
}

type TaskInsert struct {
	AssignedID       *string    `json:"assigned_id,omitempty"`
	AssignerID       *string    `json:"assigner_id,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	Desc             *string    `json:"desc,omitempty"`
	Done             *bool      `json:"done,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	ID               *int       `json:"id,omitempty"`
	MilestoneID      *int       `json:"milestone_id,omitempty"`
	TaskName         *string    `json:"task_name,omitempty"`
	TaskRef          *string    `json:"task_ref,omitempty"`
	ProjectID        *int       `json:"project_id,omitempty"`
	StoryPoints      *int       `json:"story_points,omitempty"`
	RecentProgressId *int       `json:"recent_progress"`
	FinishedDate     *time.Time `json:"finished_date,omitempty"`
}

type TaskUpdate = TaskInsert

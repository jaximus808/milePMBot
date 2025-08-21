package util

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
)

type RoleLevel int

const (
	DefaultRole RoleLevel = iota
	LeadRole
	AdminRole
	OwnerRole
)

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

func CheckDiscordPerm(userId string, guildId string, perms int64) bool {
	guild, err := discord.DiscordSession.State.Guild(guildId)
	if err != nil {
		return false
	}
	const manageGuild = discordgo.PermissionManageChannels
	const adminGuild = discordgo.PermissionAdministrator
	return perms&manageGuild != 0 || perms&adminGuild != 0 || guild.OwnerID == userId
}

type SupaDB struct{}

/**
* Supabase types
**/

// ActiveProjects
type ActiveProject struct {
	CreatedAt  time.Time `json:"created_at"`
	GuildID    *int      `json:"guildId"`
	ID         int       `json:"id"`
	PChannelID *int      `json:"pChannelId"`
	ProjectID  *int      `json:"project_id"`
}

type ActiveProjectInsert struct {
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	GuildID    *int       `json:"guildId,omitempty"`
	ID         *int       `json:"id,omitempty"`
	PChannelID *int       `json:"pChannelId,omitempty"`
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
	GuildID       *int      `json:"guild_id"`
	ID            int       `json:"id"`
	OutputChannel *int      `json:"output_channel"`
	PChannelID    *int      `json:"pchannel_id"`
	SprintEnabled *bool     `json:"sprint_enabled"`
	SprintInt     *int      `json:"sprint_int"`
	SprintMsg     *string   `json:"sprint_msg"`
	SprintN       *int      `json:"sprint_n"`
	SprintPing    *bool     `json:"sprint_ping"`
	LastPingAt    time.Time `json:"last_ping_date"`
	ProjectRef    *string   `json:"project_ref"`
}

type ProjectInsert struct {
	Completed     *bool      `json:"completed,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	CurrentMID    *int       `json:"current_mid,omitempty"`
	Desc          *string    `json:"desc,omitempty"`
	GuildID       *int       `json:"guild_id,omitempty"`
	ID            *int       `json:"id,omitempty"`
	OutputChannel *int       `json:"output_channel,omitempty"`
	PChannelID    *int       `json:"pchannel_id,omitempty"`
	SprintEnabled *bool      `json:"sprint_enabled,omitempty"`
	SprintInt     *int       `json:"sprint_int,omitempty"`
	SprintMsg     *string    `json:"sprint_msg,omitempty"`
	SprintN       *int       `json:"sprint_n,omitempty"`
	SprintPing    *bool      `json:"sprint_ping,omitempty"`
	LastPingAt    time.Time  `json:"last_ping_date,omitempty"`
	ProjectRef    *string    `json:"project_ref,omitempty"`
}

type ProjectUpdate = ProjectInsert

// Roles
type Role struct {
	CreatedAt time.Time `json:"created_at"`
	DiscordID *int      `json:"discord_id"`
	ID        int       `json:"id"`
	ProjectID *int      `json:"project_id"`
	RoleLevel int       `json:"role_level"`
}

type RoleInsert struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	DiscordID *int       `json:"discord_id,omitempty"`
	ID        *int       `json:"id,omitempty"`
	ProjectID *int       `json:"project_id,omitempty"`
	RoleLevel int        `json:"role_level,omitempty"`
}

type RoleUpdate = RoleInsert

// Tasks
type Task struct {
	AssignedID   *int       `json:"assigned_id"`
	AssignerID   *int       `json:"assigner_id"`
	CreatedAt    time.Time  `json:"created_at"`
	Desc         *string    `json:"desc"`
	Done         *bool      `json:"done"`
	Completed    *bool      `json:"completed"`
	DueDate      *time.Time `json:"due_date"`
	ID           int        `json:"id"`
	MilestoneID  *int       `json:"milestone_id,omitempty"`
	TaskName     *string    `json:"task_name"`
	TaskRef      *string    `json:"task_ref"`
	ProjectID    *int       `json:"project_id,omitempty"`
	StoryPoints  *int       `json:"story_points,omitempty"`
	FinishedDate *time.Time `json:"finished_date,omitempty"`
}

type TaskInsert struct {
	AssignedID   *int       `json:"assigned_id,omitempty"`
	AssignerID   *int       `json:"assigner_id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Desc         *string    `json:"desc,omitempty"`
	Done         *bool      `json:"done,omitempty"`
	Completed    *bool      `json:"completed,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	ID           *int       `json:"id,omitempty"`
	MilestoneID  *int       `json:"milestone_id,omitempty"`
	TaskName     *string    `json:"task_name,omitempty"`
	TaskRef      *string    `json:"task_ref,omitempty"`
	ProjectID    *int       `json:"project_id,omitempty"`
	StoryPoints  *int       `json:"story_points,omitempty"`
	FinishedDate *time.Time `json:"finished_date,omitempty"`
}

type TaskUpdate = TaskInsert

// UserAccess
type UserAccess struct {
	CreatedAt  time.Time `json:"created_at"`
	ID         int       `json:"id"`
	SupabaseId *int      `json:"user_id"`
	ProjectID  *int      `json:"project_id"`
	DiscordId  *int      `json:"discord_id"`
}

type UserAccessInsert struct {
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ID         int       `json:"id,omitempty"`
	SupabaseId *int      `json:"user_id,omitempty"`
	ProjectID  *int      `json:"project_id,omitempty"`
	DiscordId  *int      `json:"discord_id,omitempty"`
}

type UserAccessUpdate = UserAccessInsert

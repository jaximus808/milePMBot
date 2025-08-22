package projects

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	if msgInstance.GuildID == "" {
		return util.CreateHandleReport(false, output.NOT_A_CHANNEL)
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil || channel.ParentID == "" {
		return util.CreateHandleReport(false, output.NOT_A_CHANNEL)
	}

	if !util.CheckDiscordPerm(msgInstance.Member.User.ID, msgInstance.GuildID,
		msgInstance.Member.Permissions) {
		return util.CreateHandleReport(false, "❌ Missing Server Admin Permissions!")
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")
	if dateError != nil {
		return util.CreateHandleReport(false, output.FAIL_INCORRECT_DATE)
	}

	// first check if an active project exists

	_, checkActiveProject := DB.DBGetActiveProject(channel.GuildID, channel.ParentID)

	if checkActiveProject == nil {
		return util.CreateHandleReport(false, output.FAIL_ALR_PROJECT)
	}

	userId, userIdError := strconv.Atoi(msgInstance.Member.User.ID)
	channelId, channelIdError := strconv.Atoi(msgInstance.ChannelID)
	parentId, parentIdError := strconv.Atoi(channel.ParentID)
	guildId, guildIdError := strconv.Atoi(channel.GuildID)

	if userIdError != nil || channelIdError != nil || parentIdError != nil || guildIdError != nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	project, insertErr := DB.DBCreateProject(guildId, parentId, channelId, "new project!")

	if insertErr != nil || project == nil {
		util.ReportDiscordBotError(insertErr)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	// now add milestones

	milestone, msError := DB.DBCreateMilestone(project.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		DB.DBDeleteProject(project.ID)
		util.ReportDiscordBotError(msError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	userRole, roleError := DB.DBCreateRole(project.ID, userId, int(util.OwnerRole))
	if roleError != nil || userRole == nil {
		DB.DBDeleteProject(project.ID)
		util.ReportDiscordBotError(roleError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	// I NEEED TO ADD SOME TIME OF FAILURE ROLLBACK

	// everything good, now assign this project as an active project
	activeProject, activeProjctError := DB.DBCreateActiveProject(guildId, parentId, project.ID)

	if activeProjctError != nil || activeProject == nil {
		util.ReportDiscordBotError(activeProjctError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	updatedProject, updateProjectError := DB.DBUpdateCurrentMilestone(project.ID, milestone.ID)
	if updateProjectError != nil || updatedProject == nil {
		// roll back to delete incomplete project
		// roll backs will delete all previous created rows due to cascade
		DB.DBDeleteProject(project.ID)

		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	projectRefUpdate, projectRefUpdateError := DB.DBUpdateProjectRef(project.ID)
	if projectRefUpdateError != nil || projectRefUpdate == nil {
		util.ReportDiscordBotError(projectRefUpdateError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	// now begin granting access to all users who can view this message
	go util.CreateAccessRows(msgInstance, DB, channel, project.ID)

	return util.CreateHandleReport(true, output.MakeSuccessCreateProject(*projectRefUpdate.ProjectRef))
}

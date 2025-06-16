package projects

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	// TODO: MUST ADD A CHECK THAT THERE IS NO ALREADY ACTIVE PROJECT

	if msgInstance.GuildID == "" {
		return util.CreateHandleReport(false, output.NOT_A_CHANNEL)
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil || channel.ParentID == "" {
		return util.CreateHandleReport(false, output.NOT_A_CHANNEL)
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")
	if dateError != nil {
		return util.CreateHandleReport(false, output.FAIL_INCORRECT_DATE)
	}

	// first check if an active project exists

	_, checkActiveProject := util.DBGetActiveProject(channel.GuildID, channel.ParentID)

	if checkActiveProject == nil {
		return util.CreateHandleReport(false, output.FAIL_ALR_PROJECT)
	}

	userId, userIdError := strconv.Atoi(msgInstance.Member.User.ID)
	channelId, channelIdError := strconv.Atoi(msgInstance.ChannelID)
	parentId, parentIdError := strconv.Atoi(channel.ParentID)
	guildId, guildIdError := strconv.Atoi(channel.GuildID)

	if userIdError != nil || channelIdError != nil || parentIdError != nil || guildIdError != nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	project, insertErr := util.DBCreateProject(guildId, parentId, channelId, "new project!")

	if insertErr != nil || project == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	// now add milestones

	milestone, msError := util.DBCreateMilestone(project.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	userRole, roleError := util.DBCreateRole(project.ID, userId, int(util.LeadRole))
	if roleError != nil || userRole == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	//I NEEED TO ADD SOME TIME OF FAILURE ROLLBACK

	// everything good, now assign this project as an active project
	activeProject, activeProjctError := util.DBCreateActiveProject(guildId, parentId, project.ID)

	if activeProjctError != nil || activeProject == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	updatedProject, updateProjectError := util.DBUpdateCurrentMilestone(project.ID, milestone.ID)
	if updateProjectError != nil || updatedProject == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReport(true, output.SUCCESS_CREATE_PROJECT)
}

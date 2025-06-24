package projects

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ResumeProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	_, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle == nil {
		return util.CreateHandleReport(false, "❌ You can't move a project into a category with an already active project!")
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil || channel.ParentID == "" {
		return util.CreateHandleReport(false, output.NOT_A_CHANNEL)
	}

	if !util.CheckDiscordPerm(msgInstance.Member.User.ID, msgInstance.GuildID,
		msgInstance.Member.Permissions) {
		return util.CreateHandleReport(false, "❌ Missing Server Admin Permissions!")
	}

	// must be owner

	projectRef := util.GetOptionValue(args.Options, "projectref")

	selectProject, selectedProjectError := DB.DBGetProjectRef(projectRef)
	if selectedProjectError != nil || selectProject == nil {
		return util.CreateHandleReport(false, "❌ No project exists with that projectref!")
	}

	userRole, userRoleError := DB.DBGetRole(selectProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You dont have permissions to move that project here!")
	}

	if userRole.RoleLevel != int(util.OwnerRole) {
		return util.CreateHandleReport(false, "❌ You dont have permissions to move that project here!")
	}

	channelId, channelIdError := strconv.Atoi(msgInstance.ChannelID)
	parentId, parentIdError := strconv.Atoi(channel.ParentID)
	guildId, guildIdError := strconv.Atoi(channel.GuildID)

	if channelIdError != nil || parentIdError != nil || guildIdError != nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	newActiveProject, newActiveProjectError := DB.DBCreateActiveProject(guildId, parentId, selectProject.ID)

	if newActiveProjectError != nil || newActiveProject == nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	updateProject, updateProjectError := DB.DBUpdateProjectOutputChannel(selectProject.ID, channelId)
	if updateProjectError != nil || updateProject == nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReport(true, fmt.Sprintf("**🎉 Project %s has successfully resumed to this category!**\nAll tasks, milestone, and roles have been transfered!", projectRef))
}

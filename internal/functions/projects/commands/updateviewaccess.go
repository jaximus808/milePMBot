package projects

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func UpdateViewAccess(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	userRole, userRoleError := DB.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	if userRole.RoleLevel < int(util.AdminRole) {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)
	if err != nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}
	go func(msgInstance *discordgo.InteractionCreate, DB util.DBClient, channel *discordgo.Channel, projectID int) {
		report := util.CreateAccessRows(msgInstance, DB, channel, projectID)
		if !report.GetSuccess() {
			util.ReportDiscordBotReport(report)
		}
	}(msgInstance, DB, channel, currentProject.ID)

	return util.CreateHandleReport(true, "Updating View Access for this Discord Project for Web View! It may take a few seconds to update.")
}

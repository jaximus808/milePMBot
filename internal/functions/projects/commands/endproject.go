package projects

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func EndProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	// must be owner
	userRole, userRoleError := util.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	if userRole.RoleLevel != int(util.OwnerRole) {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	removeError := util.DBEndActiveProject(currentProject.ID)

	if removeError != nil {
		return util.CreateHandleReport(false, "❌ Something went wrong making project inactive")
	}

	return util.CreateHandleReport(true, fmt.Sprintf("**🗑️ Project %s has been made inactive**\nThe project's data won't be deleted and can be resumed in any other category in any server MilePM is present", *currentProject.ProjectRef))

}

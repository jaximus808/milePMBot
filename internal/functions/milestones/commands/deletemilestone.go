package milestones

import (
	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func DeleteMilestone(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}
	userRole, userRoleError := util.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	if userRole.RoleLevel < int(util.AdminRole) {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}
	msRef := util.GetOptionValue(args.Options, "msref")

	milestone, milestoneError := util.DBGetMilestoneWithRef(msRef, currentProject.ID)

	if milestoneError != nil || milestone == nil {
		return util.CreateHandleReport(false, "❌ Could not find milestone, double check the milestone!")
	}

	if *currentProject.CurrentMID == milestone.ID {
		return util.CreateHandleReport(false, "❌ You can't delete an active milestone!")
	}
	deleteErr := util.DBDeleteMilestone(milestone.ID)
	if deleteErr != nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReport(true, "Successfully deleted milestone!")
}

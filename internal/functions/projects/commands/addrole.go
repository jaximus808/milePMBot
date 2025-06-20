package projects

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func AddRole(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	op := util.GetOptionValue(args.Options, "op")
	user := util.GetOptionValue(args.Options, "user")
	role := util.GetOptionValue(args.Options, "role")
	re := regexp.MustCompile(`<@!?(\d+)>`)
	match := re.FindStringSubmatch(user)
	if len(match) != 2 {
		return util.CreateHandleReport(false, "You need to be @ a user!")
	}
	userId := match[1]

	userRole, userRoleError := util.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	userIdNum, idError := strconv.Atoi(userId)
	if idError != nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	// good to go now

	if op == "add" {
		roleInt := 0

		if role == "admin" {
			roleInt = int(util.AdminRole)
		} else if role == "lead" {
			roleInt = int(util.LeadRole)
		}

		if userRole.RoleLevel <= roleInt {
			return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
		}

		roleExist, roleCheckError := util.DBGetRole(currentProject.ID, userId)
		if roleCheckError == nil || roleExist != nil {
			return util.CreateHandleReport(false, "❌ A user can't have two roles!")
		}
		userRole, roleError := util.DBCreateRole(currentProject.ID, userIdNum, int(roleInt))
		if roleError != nil || userRole == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}
		return util.CreateHandleReport(true, fmt.Sprintf("Success! <@%s> has been given the role: **%s**", userId, role))
	} else if op == "remove" {
		roleExist, roleCheckError := util.DBGetRole(currentProject.ID, userId)
		if roleCheckError != nil || roleExist == nil {
			return util.CreateHandleReport(false, "❌ This user doesn't have a role!")
		}

		roleInt := roleExist.RoleLevel
		if userRole.RoleLevel <= roleInt {
			return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
		}

		// do later
		deleteError := util.DBDeleteRole(roleExist.ID)
		if deleteError != nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReport(true, fmt.Sprintf("<@%s> has been removed of their role", userId))
	}

	return util.CreateHandleReport(false, "❌ Expected [operation] [user] [role]")
}

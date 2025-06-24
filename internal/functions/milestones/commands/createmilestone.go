package milestones

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateMilestone(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}
	userRole, userRoleError := DB.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	if userRole.RoleLevel < int(util.AdminRole) {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")

	if dateError != nil {
		return util.CreateHandleReport(false, output.FAIL_INCORRECT_DATE)
	}

	milestoneExist := DB.DBMilestoneExistDate(currentProject.ID, &msDate)

	if milestoneExist {
		return util.CreateHandleReport(false, output.FAIL_MS_SAME_DATE)
	}

	// now add milestones

	milestone, msError := DB.DBCreateMilestone(currentProject.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		util.ReportDiscordBotError(msError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	// create Ref
	milestoneRef, refError := DB.DBUpdateMilestoneRef(milestone.ID, *currentProject.ProjectRef)
	if refError != nil || milestoneRef == nil {
		util.ReportDiscordBotError(refError)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}
	return util.CreateHandleReportAndOutput(true,
		"successfully created milestone with milestoneref : "+*milestoneRef.MilestoneRef, &discordgo.MessageEmbed{
			Title:       "🪜 New Milestone Created",
			Description: fmt.Sprintf("A new milestone **%s** has been added to the project!", *milestone.DisplayName),
			Color:       0x5865F2, // Discord blurple
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Milestone Name", Value: msName, Inline: false},
				{Name: "Milestone Ref", Value: *milestoneRef.MilestoneRef, Inline: false},
				{Name: "Due Date", Value: msDate.Format("January 2, 2006"), Inline: false}, // if available
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}, strconv.Itoa(*currentProject.OutputChannel),
	)
}

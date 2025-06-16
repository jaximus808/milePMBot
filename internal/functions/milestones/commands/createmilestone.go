package milestones

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateMilestone(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")
	if dateError != nil {
		return util.CreateHandleReport(false, output.FAIL_INCORRECT_DATE)
	}

	milestoneExist := util.DBMilestoneExistDate(currentProject.ID, &msDate)

	if milestoneExist {
		return util.CreateHandleReport(false, output.FAIL_MS_SAME_DATE)
	}

	// now add milestones

	milestone, msError := util.DBCreateMilestone(currentProject.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	return util.CreateHandleReportAndOutput(true,
		"successfully created milestone with id: "+strconv.Itoa(milestone.ID), &discordgo.MessageEmbed{
			Title:       "🪜 New Milestone Created",
			Description: fmt.Sprintf("A new milestone **%s** has been added to the project!", *milestone.DisplayName),
			Color:       0x5865F2, // Discord blurple
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Milestone Name", Value: msName, Inline: false},
				{Name: "Due Date", Value: msDate.Format("January 2, 2006"), Inline: false}, // if available
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}, strconv.Itoa(*currentProject.OutputChannel),
	)
}

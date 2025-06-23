package milestones

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func ListMilestones(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance, DB)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	milestoneList, milestoneListError := DB.DBGetMilestoneListDescending(currentProject.ID)

	if milestoneListError != nil || milestoneList == nil {
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	emeddedMessage := &discordgo.MessageEmbed{
		Title:       "🧭 Milestone Map",
		Description: "For: Current Project",
		Color:       0x3498DB, // Green
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	milestoneMap := util.ParseMilestoneList(milestoneList, *currentProject.CurrentMID)

	currentField := &discordgo.MessageEmbedField{
		Name:   "🚀 Current",
		Value:  milestoneMap.CurrentMilestone,
		Inline: false,
	}

	upcomingField := &discordgo.MessageEmbedField{
		Name:   "📋 Upcoming",
		Inline: true,
	}

	previousField := &discordgo.MessageEmbedField{
		Name:   "📨 Previous",
		Inline: true,
	}

	if len(milestoneMap.Upcoming) > 0 {
		upcomingField.Value = strings.Join(milestoneMap.Upcoming, "\n")
	} else {
		upcomingField.Value = "No Upcoming Milestones"
	}

	if len(milestoneMap.Previous) > 0 {
		previousField.Value = strings.Join(milestoneMap.Previous, "\n")
	} else {
		previousField.Value = "No Previous Milestones"
	}

	emeddedMessage.Fields = []*discordgo.MessageEmbedField{
		currentField,
		upcomingField,
		previousField,
	}
	return util.CreateHandleReportAndOutput(
		true,
		"Map Made!",
		emeddedMessage,
		msgInstance.ChannelID,
	)

}

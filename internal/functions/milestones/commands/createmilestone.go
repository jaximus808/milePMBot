package milestones

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateMilestone(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")
	if dateError != nil {
		return util.CreateHandleReport(false, "incorrect date format, expect MM/DD/YYYY")
	}

	milestoneExist := util.DBMilestoneExistDate(currentProject.ID, &msDate)

	if milestoneExist {
		return util.CreateHandleReport(false, "two milestones can't have the same date!!!")
	}

	// now add milestones

	milestone, msError := util.DBCreateMilestone(currentProject.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		return util.CreateHandleReport(false, "failed to make milestone tied to project")
	}

	return util.CreateHandleReport(true, "successfully created milestone with id: "+strconv.Itoa(milestone.ID))
}

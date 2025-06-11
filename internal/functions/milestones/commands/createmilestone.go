package milestones

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateMilestone(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {

	if len(args) < 3 {
		return util.CreateHandleReport(false, "expected 3 args [ms_name] [date] [desc]")
	}

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, "failed to get active project")
	}

	msName := args[0]
	msDate, dateError := time.Parse("01/02/2006", args[1])
	msDesc := strings.Join(args[2:], " ")
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

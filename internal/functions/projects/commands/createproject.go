package projects

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateProject(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport {

	if msgInstance.GuildID == "" {
		return util.CreateHandleReport(false, "Not in a discord server")
	}

	if len(args) < 3 {
		return util.CreateHandleReport(false, "could not parse instructoons")
	}
	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil {
		return util.CreateHandleReport(false, "channel failed!")
	}

	if channel.Type != 4 {
		return util.CreateHandleReport(false, "message must be in a category!")
	}

	project, insertErr := util.DBCreateProject(channel.GuildID, channel.ParentID, msgInstance.ChannelID)

	if insertErr != nil || project == nil {
		return util.CreateHandleReport(false, "failed to make project")
	}

	msName := args[0]
	msDate, dateError := time.Parse("01/02/2006", args[1])
	msDesc := strings.Join(args[2:], " ")
	if dateError != nil {
		return util.CreateHandleReport(false, "incorrect date format, expect MM/DD/YYYY")
	}

	// now add milestones

	milestone, msError := util.DBCreateMilestone(project.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		return util.CreateHandleReport(false, "failed to make milestone tied to project")
	}

	userRole, roleError := util.DBCreateRole(project.ID, msgInstance.Author.ID, int(util.LeadRole))
	if roleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "failed to make user role tied to project")
	}

	return util.CreateHandleReport(true, "Successfully added project!")
}

package projects

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func CreateProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	// TODO: MUST ADD A CHECK THAT THERE IS NO ALREADY ACTIVE PROJECT

	if msgInstance.GuildID == "" {
		return util.CreateHandleReport(false, "Not in a discord server")
	}

	channel, err := discord.DiscordSession.Channel(msgInstance.ChannelID)

	if err != nil {
		return util.CreateHandleReport(false, "channel failed!")
	}

	if channel.ParentID == "" {
		return util.CreateHandleReport(false, "message must be in a category!")
	}

	msName := util.GetOptionValue(args.Options, "msname")
	msDate, dateError := time.Parse("01/02/2006", util.GetOptionValue(args.Options, "msdate"))
	msDesc := util.GetOptionValue(args.Options, "desc")
	if dateError != nil {
		return util.CreateHandleReport(false, "incorrect date format, expect MM/DD/YYYY: "+dateError.Error())
	}

	// first check if an active project exists

	_, checkActiveProject := util.DBGetActiveProject(channel.GuildID, channel.ParentID)

	if checkActiveProject == nil {
		return util.CreateHandleReport(false, "There already is an active project for this category!")
	}

	log.Printf("maowgood?")
	project, insertErr := util.DBCreateProject(channel.GuildID, channel.ParentID, msgInstance.ChannelID, "new project!")

	if insertErr != nil || project == nil {
		return util.CreateHandleReport(false, "failed to make project")
	}

	// now add milestones

	milestone, msError := util.DBCreateMilestone(project.ID, msName, &msDate, msDesc)
	if msError != nil || milestone == nil {
		return util.CreateHandleReport(false, "failed to make milestone tied to project")
	}

	userRole, roleError := util.DBCreateRole(project.ID, msgInstance.Member.User.ID, int(util.LeadRole))
	if roleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "failed to make user role tied to project")
	}
	//I NEEED TO ADD SOME TIME OF FAILURE ROLLBACK

	// everything good, now assign this project as an active project
	activeProject, activeProjctError := util.DBCreateActiveProject(channel.GuildID, channel.ParentID, project.ID)

	if activeProjctError != nil || activeProject == nil {
		return util.CreateHandleReport(false, "failed to create an active project :(")
	}
	updatedProject, updateProjectError := util.DBUpdateCurrentMilestone(project.ID, milestone.ID)
	if updateProjectError != nil || updatedProject == nil {
		return util.CreateHandleReport(false, "failed to create an update milestone  :(")
	}

	return util.CreateHandleReport(true, "Successfully added project!")
}

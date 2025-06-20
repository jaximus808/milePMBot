package projects

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func SettingProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil || currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}

	setting := util.GetOptionValue(args.Options, "setting")
	value := util.GetOptionValue(args.Options, "value")

	userRole, userRoleError := util.DBGetRole(currentProject.ID, msgInstance.Member.User.ID)

	if userRoleError != nil || userRole == nil {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	if userRole.RoleLevel < int(util.AdminRole) {
		return util.CreateHandleReport(false, "❌ You don't have the valid permission for this command")
	}

	switch setting {
	case "output":
		re := regexp.MustCompile(`<#!?(\d+)>`)
		match := re.FindStringSubmatch(value)
		if len(match) != 2 {
			return util.CreateHandleReport(false, "❌ You need to give a valid channel")
		}
		outChannel := match[1]

		outChannelId, errOut := strconv.Atoi(outChannel)
		if errOut != nil {
			util.CreateHandleReport(false, "❌ You need to give a valid channel")
		}

		//need a check to make sure the channel is within the cateogry

		newChannel, channelErr := discord.DiscordSession.Channel(outChannel)

		parentId, praseError := strconv.Atoi(newChannel.ParentID)

		if channelErr != nil || praseError != nil {
			return util.CreateHandleReport(false, "❌ You need to give a valid channel")
		}

		if parentId != *currentProject.PChannelID {
			return util.CreateHandleReport(false, "❌ You need to give a channel within the same category")
		}

		updatedProject, errUpdate := util.DBUpdateProjectOutputChannel(currentProject.ID, outChannelId)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			"✅ Output channel updated sucessfull!",
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: "## Project Output Channel Has Been Updated to Here",
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	case "description":
		updatedProject, errUpdate := util.DBUpdateProjectDescription(currentProject.ID, value)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			"✅ Project Description Updated!",
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: fmt.Sprintf("## Project Description to:\n%s", value),
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	case "message":
		updatedProject, errUpdate := util.DBUpdateProjectSprintDesc(currentProject.ID, value)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			"✅ Project Sprint Message Updated!",
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: fmt.Sprintf("## Project Description to:\n%s", value),
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	case "sprints":
		value = strings.ToLower(value)
		var toggledSprint bool
		var stringSpring string
		if value == "y" {
			toggledSprint = true
			stringSpring = "Sprints Enabled!"
		} else if value == "n" {
			toggledSprint = false
			stringSpring = "Sprints Disabled"
		} else {
			return util.CreateHandleReport(false, "❌ expected Y or N")
		}
		updatedProject, errUpdate := util.DBUpdateProjectSprints(currentProject.ID, toggledSprint)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			fmt.Sprintf("✅ Project %s", stringSpring),
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: fmt.Sprintf("## Project %s", stringSpring),
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	case "duration":
		sprintDuration, sprintDurationError := strconv.Atoi(value)
		if sprintDurationError != nil {
			return util.CreateHandleReport(false, "❌ Expected a Number for the Sprint Duration. Ex: 3 -> 3 weeks")
		}

		updatedProject, errUpdate := util.DBUpdateProjectSprintDuration(currentProject.ID, sprintDuration)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			"✅ Project Sprint Duration Updated",
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: fmt.Sprintf("# Project Sprint Duration Updated To: %d", sprintDuration),
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	case "pings":
		value = strings.ToLower(value)
		var toggledSprint bool
		var stringSpring string
		if value == "y" {
			toggledSprint = true
			stringSpring = "Pings Enabled!"
		} else if value == "n" {
			toggledSprint = false
			stringSpring = "Pings Disabled"
		} else {
			return util.CreateHandleReport(false, "❌ expected Y or N")
		}
		updatedProject, errUpdate := util.DBUpdateProjectPings(currentProject.ID, toggledSprint)

		if errUpdate != nil || updatedProject == nil {
			return util.CreateHandleReport(false, output.FAILURE_SERVER)
		}

		return util.CreateHandleReportAndOutput(
			true,
			fmt.Sprintf("✅ Project %s", stringSpring),
			&discordgo.MessageEmbed{
				Title:       "🔁 Project Updated",
				Description: fmt.Sprintf("## Project %s", stringSpring),
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)
	}

	return util.CreateHandleReport(false, "❌ Unexpected Option")
}

package projects

import (
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func SettingProject(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport {

	currentProject, errorHandle := util.SetUpProjectInfo(msgInstance)

	if errorHandle != nil {
		return errorHandle
	}

	if currentProject == nil {
		return util.CreateHandleReport(false, output.NO_ACTIVE_PROJECT)
	}
	setting := util.GetOptionValue(args.Options, "setting")
	value := util.GetOptionValue(args.Options, "value")

	switch setting {
	case "output":
		re := regexp.MustCompile(`<#([0-9])>`)
		match := re.FindStringSubmatch(value)
		if len(match) != 2 {
			return util.CreateHandleReport(false, "You need to give a channel")
		}
		outChannel := match[1]

		outChannelId, errOut := strconv.Atoi(outChannel)
		if errOut != nil {
			util.CreateHandleReport(false, "expected output id")
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
				Description: "Project Output Channel Has Been Updated to Here",
				Color:       0x3498DB, // Orange
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			strconv.Itoa(*updatedProject.OutputChannel),
		)

	}

	return util.CreateHandleReport(false, "❌ Unexpected Option")
}

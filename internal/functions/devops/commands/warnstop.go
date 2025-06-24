package devops

import (
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	output "github.com/jaximus808/milePMBot/internal/ouput/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func WarnDeployments(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	// fuck i gotta do some bs here

	if msgInstance.GuildID != os.Getenv("ADMIN_GUILD") || msgInstance.Member.User.ID != os.Getenv("ADMIN_USER") {
		return util.CreateHandleReport(false, "Didn't recongnize your command")
	}

	projects, err := DB.DBGetAllPingProjects()

	if err != nil || projects == nil {
		util.ReportDiscordBotError(err)
		return util.CreateHandleReport(false, output.FAILURE_SERVER)
	}

	emeddMessage := &discordgo.MessageEmbed{
		Title:       "⚙️ Bot Maintainence Scheduled",
		Description: "MilestonePM will be paused for maintainence and updates in around 5 minutes\nWe apologize for the inconvience 😅",
	}
	for _, project := range *projects {
		discord.DiscordSession.ChannelMessageSendEmbed(strconv.Itoa(*project.OutputChannel), emeddMessage)
	}

	return util.CreateHandleReport(true, "Message sent successfully")

}

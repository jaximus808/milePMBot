package jobs

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/util"
)

func WeeklyRemindProjects() {
	projects, err := util.DBGetAllPingProjects()

	if err != nil {
		log.Print(err.Error())
		return
	}

	log.Printf("beginning weekly remind")
	for _, project := range *projects {
		// this might not scale well
		if isAtLeastNWeeksAgo(project.LastPingAt, *project.SprintN) {
			sendWeeklyMessage(&project)
			util.DBUpdateResetSprintDuration(project.ID)
		}
	}

	log.Printf("finished weekly remind")
}
func isAtLeastNWeeksAgo(t time.Time, n int) bool {
	// time.Since(t) is the same as time.Now().Sub(t)
	return time.Since(t) >= time.Duration(n)*7*24*time.Hour
}

func sendWeeklyMessage(project *util.Project) {

	tasksMs, err := util.DBGetInProgressAndCompetedTask(project.ID, *project.CurrentMID)

	if err != nil {
		log.Print(err.Error())
		return
	}
	var taskReport *util.TaskReport
	if *project.SprintPing {
		taskReport = util.ParseTaskListWeeklyWithPing(tasksMs)
	} else {
		taskReport = util.ParseTaskListWeekly(tasksMs, strconv.Itoa(*project.GuildID))
	}

	emeddMessage := &discordgo.MessageEmbed{
		Title:       "🗓️ Sprint Update",
		Description: *project.SprintMsg,
		Color:       0x7289DA, // Discord blurple
		Timestamp:   time.Now().Format(time.RFC3339),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🔍 In Review",
				Value:  strings.Join(taskReport.InProgress, "\n"),
				Inline: false,
			},
			{
				Name:   "🚧 In Progress",
				Value:  strings.Join(taskReport.InProgress, "\n"),
				Inline: false,
			},
		},
	}

	discord.DiscordSession.ChannelMessageSendEmbed(strconv.Itoa(*project.OutputChannel), emeddMessage)
}

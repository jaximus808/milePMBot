package projects

import (
	"github.com/bwmarrin/discordgo"
	projects "github.com/jaximus808/milePMBot/internal/functions/projects/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport){
	"start":  projects.CreateProject,
	"set":    projects.SettingProject,
	"role":   projects.AddRole,
	"move":   projects.MoveProject,
	"end":    projects.EndProject,
	"resume": projects.ResumeProject,
	"info":   projects.ProjectInfo,
}

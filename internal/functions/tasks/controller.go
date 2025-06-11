package tasks

import (
	"github.com/bwmarrin/discordgo"
	tasks "github.com/jaximus808/milePMBot/internal/functions/tasks/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport){
	"create":   tasks.AddTask,
	"assign":   tasks.AssignTask,
	"complete": tasks.CompleteTask,
	"approve":  tasks.ApproveTask,
	"reject":   tasks.RejectTask,
}

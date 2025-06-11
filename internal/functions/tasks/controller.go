package tasks

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.MessageCreate, args []string) *util.HandleReport){}

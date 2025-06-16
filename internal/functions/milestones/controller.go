package milestones

import (
	"github.com/bwmarrin/discordgo"
	milestones "github.com/jaximus808/milePMBot/internal/functions/milestones/commands"
	"github.com/jaximus808/milePMBot/internal/util"
)

var CommandMap = map[string](func(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption) *util.HandleReport){
	"create": milestones.CreateMilestone,
	"move":   milestones.MoveMilestone,
	"map":    milestones.ListMilestones,
}

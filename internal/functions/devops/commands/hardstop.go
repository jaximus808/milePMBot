package devops

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jaximus808/milePMBot/internal/util"
)

func HardStop(msgInstance *discordgo.InteractionCreate, args *discordgo.ApplicationCommandInteractionDataOption, DB util.DBClient) *util.HandleReport {
	if msgInstance.GuildID != os.Getenv("ADMIN_GUILD") || msgInstance.Member.User.ID != os.Getenv("ADMIN_USER") {
		return util.CreateHandleReport(false, "Didn't recongnize your command")
	}

	log.Printf("recieved forced shutdown command, stopping service")

	os.Exit(0)

	return nil
}

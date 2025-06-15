package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var DiscordSession *discordgo.Session

func InitalizeDiscordGo() error {

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))

	if err != nil {
		return err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuilds | discordgo.IntentsGuildMembers
	DiscordSession = session
	return nil
}

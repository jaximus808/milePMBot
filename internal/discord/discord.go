package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func InitalizeDiscordGo() (*discordgo.Session, error) {

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))

	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent | discordgo.IntentsGuilds
	return session, nil
}

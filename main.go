package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jaximus808/milePMBot/internal/discord"
	integration "github.com/jaximus808/milePMBot/internal/integration/discord"
	"github.com/jaximus808/milePMBot/internal/jobs"
	"github.com/jaximus808/milePMBot/internal/supabaseutil"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print(err.Error())
	}

	supabaseErr := supabaseutil.InitializeSupabase()

	if supabaseErr != nil {
		fmt.Print(supabaseErr.Error())
		return
	}

	discErr := discord.InitalizeDiscordGo()

	if discErr != nil {
		fmt.Print(discErr.Error())
	}

	discord.DiscordSession.AddHandler(integration.MainHandler)

	openErr := discord.DiscordSession.Open()
	if openErr != nil {
		log.Fatalf("Cannot open the session: %v", openErr)
	}
	defer discord.DiscordSession.Close()

	integration.ClearCommands(discord.DiscordSession, os.Getenv("DEPLOY_GUILD"))
	integration.RegisterCommands(discord.DiscordSession, os.Getenv("DEPLOY_GUILD"), os.Getenv("ADMIN_GUILD"))

	// start the cron job
	s, err := jobs.StartSprintUpdateJob()
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer s.Shutdown()

	s.Start()

	for _, job := range s.Jobs() {
		log.Printf("running job %s", job.ID().String())
	}

	log.Println("Bot and Job Online")
	discord.DiscordSession.ChannelMessageSend(os.Getenv("OUTPUT_LOG_CHANNEL"), "🚀 MilestonePM bot is online")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
	// now start adding handler
}

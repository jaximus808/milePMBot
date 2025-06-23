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
	envErr := godotenv.Load(".env")
	if envErr != nil {
		fmt.Printf("No env file to be found")
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

	integration.ClearCommands(discord.DiscordSession, "738509536520044575")
	integration.RegisterCommands(discord.DiscordSession, "738509536520044575")

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
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
	// now start adding handler

}

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jaximus808/milePMBot/internal/discord"
	"github.com/jaximus808/milePMBot/internal/functions"
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

	discordSession, discErr := discord.InitalizeDiscordGo()

	if discErr != nil {
		fmt.Print(discErr.Error())
	}
	discordSession.AddHandler(functions.MainHandler)

	openErr := discordSession.Open()
	if openErr != nil {
		log.Fatalf("Cannot open the session: %v", openErr)
	}
	defer discordSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
	// now start adding handler

}

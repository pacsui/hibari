package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	log "github.com/charmbracelet/log"
	"github.com/pacsui/threadsinchannel/handlers"
)

// Flags
var (
	Debug bool
)

func init() {
	flag.BoolVar(&Debug, "d", false, "Set debug mode")
	flag.Parse()
	log.SetReportCaller(Debug)
	if Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in Debug!")
	}
	dConVal, err := handlers.ReadConfigFile("config.yaml")
	if err != nil {
		log.Error(err)
		return
	}
	handlers.DiscordBotConfigValues = dConVal
}

func main() {
	s, _ := discordgo.New("Bot " + handlers.DiscordBotConfigValues.DiscordConfig.Auth.Token)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot running")
	})

	//TODO: add []interface{}
	s.AddHandler(handlers.HandleStarBoardAdd)
	s.AddHandler(handlers.HandleMessageInThreads)
	s.AddHandler(handlers.HandleStarBoardDel)

	done := make(chan struct{})
	go handlers.PollingServiceToCrossPost(done, s)

	// go func(s *discordgo.Session) {
	// 	for i := range handlers.QUEUE {
	// 		log.Info(i)
	// 	}
	// }(s)

	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Info("Graceful shutdown")

}

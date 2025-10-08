package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	log "github.com/charmbracelet/log"
	"github.com/pacsui/threadsinchannel/handlers"
	yaml "gopkg.in/yaml.v3"
)

// Flags
var (
	DiscordConfigValues handlers.BotConfig
	Debug               bool
)

func init() {
	flag.BoolVar(&Debug, "d", false, "Set debug mode")
	flag.Parse()
	log.SetReportCaller(Debug)
	if Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in Debug!")
	}

	configBytes, err := os.ReadFile("config.yaml")
	log.Debug(string(configBytes))
	if err != nil {
		log.Warn("Config file not found config.yaml")
		fl, err := os.Create("config.yaml")
		if err != nil {
			log.Errorf("unable to create file : %s", err.Error())
		}
		emptyBotConfig := handlers.BotConfig{}
		yamlDefault, err := yaml.Marshal(emptyBotConfig)

		_, err = fl.WriteString(string(yamlDefault))
		if err != nil {
			log.Warn("Unable to create default config file : %s", err.Error())
		}
	}
	err = yaml.Unmarshal(configBytes, &DiscordConfigValues)
	if err != nil {
		log.Error("unable to unmarshal!? : %s", err.Error())
		return
	}
	log.Debug(DiscordConfigValues)

}

func main() {
	s, _ := discordgo.New("Bot " + DiscordConfigValues.Auth.Token)
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

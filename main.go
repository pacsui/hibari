package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	log "github.com/charmbracelet/log"
	"github.com/pacsui/threadsinchannel/handlers"
)

// Flags
var (
	Debug       bool
	HandlerList []handlers.Handler
)

func init() {
	flag.BoolVar(&Debug, "d", false, "Set debug mode")
	flag.Parse()
	log.SetReportCaller(Debug)
	if Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in Debug!")
		log.Debugf("main PPID: %d", os.Getpid())
	}
	dConVal, err := handlers.ReadConfigFile("config.yaml")
	if err != nil {
		log.Error(err)
		return
	}
	handlers.DiscordBotConfigValues = dConVal

	HandlerList = []handlers.Handler{
		handlers.Handler{
			Name:     "starboard handler",
			Function: handlers.HandleStarBoardAdd,
		},
		handlers.Handler{
			Name:     "msg thread handler",
			Function: handlers.HandleMessageInThreads,
		},
		handlers.Handler{
			Name:     "starboard del",
			Function: handlers.HandleStarBoardDel,
		},
		// handlers.Handler{
		// 	Name:     "avatar handler",
		// 	Function: handlers.HandleAvatarEmbedReply,
		// },
		// handlers.Handler{
		// 	Name:     "impersonation handler",
		// 	Function: handlers.HandleImpersonation,
		// },
	}

}

func main() {
	s, _ := discordgo.New("Bot " + handlers.DiscordBotConfigValues.DiscordConfig.Auth.Token)
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Bot running as %s", s.State.User.DisplayName())
	})

	for _, handler := range HandlerList {
		s.AddHandler(handler.Function)
		log.Infof("Added Handler : %s", handler.Name)
	}
	done := make(chan struct{})
	go handlers.PollingServiceToCrossPost(done, s)
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

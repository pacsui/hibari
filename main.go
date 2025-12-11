package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

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

	r, _ := InitRedis()

	HandlerList = []handlers.Handler{
		{
			Name:     "starboard_handler",
			Function: handlers.HandleStarBoardAdd,
			File:     "starboard.go",
		},
		{
			Name:     "thread_creator",
			Function: handlers.HandleMessageInChannelPool,
			File:     "channelthread.go",
		},
		{
			Name:     "old commands handler",
			Function: handlers.OnMessageOldCommandHandler,
			File:     "miscported.go",
		},
		{
			Name: "cap bg counter",
			Function: func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
				handlers.CapBoardHandler(s, m, r)
			},
			File: "capboard.go",
		},
		{
			Name: "cap commands",
			Function: func(s *discordgo.Session, m *discordgo.MessageCreate) {
				handlers.CapBoardCommandHandler(s, m, r)
			},
			File: "capboard.go",
		},
		{
			Name:     "confession handler",
			Function: handlers.ConfessionMessageHandler,
			File:     "confession.go",
		},
	}

}

func main() {
	s, _ := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Infof("Bot running as %s", s.State.User.DisplayName())
	})

	for _, handler := range HandlerList {
		s.AddHandler(handler.Function)
		log.Infof("Added Handler : %s", handler.Name)
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if strings.HasPrefix(m.Content, handlers.C("mixins")) {
			mixinList := "Enabled Mixins :\n```"
			for _, handler := range HandlerList {
				mixinList += fmt.Sprintf("- %s (%s)\n", handler.Name, handler.File)
			}
			mixinList += "\n```"
			s.ChannelMessageSend(m.ChannelID, mixinList)
		}
	})
	// done := make(chan struct{})
	// go handlers.PollingServiceToCrossPost(done, s)
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

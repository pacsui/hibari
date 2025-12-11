package handlers

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var threadRequests = []string{
	"%s, we recommend using threads for new topics.",
}

type ChanMsgKeyType struct {
	ChanID string
	MsgID  string
}

func ChanMsgKey(cID string, mID string) ChanMsgKeyType {
	if len(mID) > 0 && len(cID) > 0 {
		return ChanMsgKeyType{
			ChanID: cID,
			MsgID:  mID,
		}
	}
	return ChanMsgKeyType{}
}

func Hlog(s *discordgo.Session, content string) {
	// TODO: Implement discord channel logging?
	s.ChannelMessageSend("1427259484299857961", content)
	log.Info(content)
}

func GetRandomThreadRequest(username string) string {
	if len(threadRequests) == 0 {
		return "Error: No thread requests available."
	}
	randomIndex := rand.Intn(len(threadRequests))
	return fmt.Sprintf(threadRequests[randomIndex], username)
}

func ReadConfigFile(filepath string) (BotConfig, error) {
	// Reads `config.yaml` and populates BotConfig struct
	BConfig := BotConfig{}
	configBytes, err := os.ReadFile(filepath)
	log.Debug(string(configBytes))
	if err != nil {
		log.Warn("Config file not found config.yaml")
		fl, err := os.Create("config.yaml")
		if err != nil {
			log.Errorf("unable to create file : %s", err.Error())
		}
		emptyBotConfig := BotConfig{}
		yamlDefault, err := yaml.Marshal(emptyBotConfig)
		if err != nil {
			log.Warn("unable to write bot empty config?", err.Error())
		}

		_, err = fl.WriteString(string(yamlDefault))
		if err != nil {
			log.Warn("Unable to create default config file : %s", err.Error())
		}
	}
	err = yaml.Unmarshal(configBytes, &BConfig)
	if err != nil {
		log.Error("unable to unmarshal!? : %s", err.Error())
		return BotConfig{}, err

	}
	return BConfig, nil
}

func C(cmd_str string) string {
	// convert `command` to `>>command`
	return DiscordBotConfigValues.CommandPrefix + cmd_str
}

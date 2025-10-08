package handlers

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
)

var threadRequests = []string{
	"%s, we recommend using threads for new topics.",
}

func GetRandomThreadRequest(username string) string {
	if len(threadRequests) == 0 {
		return "Error: No thread requests available."
	}
	randomIndex := rand.Intn(len(threadRequests))
	return fmt.Sprintf(threadRequests[randomIndex], username)
}

func ReadConfigFile(filepath string) (BotConfig, error) {
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

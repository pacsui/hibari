package handlers

import (
	"fmt"
	"math/rand"
)

type AuthConfig struct {
	Token string `yaml:"token"`
}

type ChannelsConfig struct {
	PictureChannel string `yaml:"picture-channel"`
}

type DiscordConfig struct {
	Auth     AuthConfig     `yaml:"auth"`
	Channels ChannelsConfig `yaml:"channels"`
}

type RedisConfig struct {
	Endpoint string `yaml:"endpoint"`
	Password string `yaml:"password"`
	DbNum    int    `yaml:"dbnum"`
	DoneVal  int    `yaml:"doneval"`
}

type BotConfig struct {
	DiscordConfig `yaml:"discord"`
	RedisConfig   `yaml:"redis"`
}

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

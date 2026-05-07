package handlers

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

var threadRequests = []string{
	"%s, we recommend using threads for new topics.",
}

type ChanMsgKeyType struct {
	ChanID string
	MsgID  string
}

func GenerateColorHex(uid string) int {
	salt := time.Now().Unix() / 3600
	input := fmt.Sprintf("%s-%d", uid, salt)
	hash := sha256.Sum256([]byte(input))

	h := float64(hash[0]) / 255.0
	s := 0.70 + (float64(hash[1])/255.0)*0.25
	l := 0.80 + (float64(hash[2])/255.0)*0.10

	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		hueToRGB := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	red := int((r * 255.0) + 0.5)
	green := int((g * 255.0) + 0.5)
	blue := int((b * 255.0) + 0.5)

	return (red << 16) | (green << 8) | blue
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
	spew.Dump(BConfig)
	return BConfig, nil
}

func MessageURL(channelID string, messageID string) string {
	return "https://discord.com/channels/" + DiscordBotConfigValues.DiscordConfig.GuildID + "/" + channelID + "/" + messageID

}

func C(cmd_str string) string {
	// convert `command` to `>>command`
	return DiscordBotConfigValues.CommandPrefix + cmd_str
}

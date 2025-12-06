package handlers

var (
	DiscordBotConfigValues BotConfig
)

type AuthConfig struct {
	Token string `yaml:"token"`
}

type StarConfig struct {
	Threshold int `yaml:"threshold"`
}

type SettingsConfig struct {
	Star StarConfig `yaml:"star"`
}

type DiscordConfig struct {
	CommandPrefix             string      `yaml:"cmd_prefix"`
	GuildIDs                  []string    `yaml:"guildIDs"`
	Auth                      AuthConfig  `yaml:"auth"`
	Channels                  []string    `yaml:"threads-filtered"`
	AllowedURLPrefix          []string    `yaml:"allowed-urls-threads-filtered"`
	StarBoardChannel          string      `yaml:"starboard-channel"`
	StarBoardFilteredChannels []string    `yaml:"starboard-disabled-channels"`
	Redis                     RedisConfig `yaml:"redis"`
	HelloReply                []string    `yaml:"hello-replies"`
}

type RedisConfig struct {
	Endpoint string `yaml:"endpoint"`
	Password string `yaml:"password"`
	DbNum    int    `yaml:"dbnum"`
	DoneVal  int    `yaml:"doneval"`
}

type BotConfig struct {
	SettingsConfig `yaml:"settings"`
	DiscordConfig  `yaml:"discord"`
	RedisConfig    `yaml:"redis"`
}

type Handler struct {
	Name     string
	Function interface{}
	File     string
}

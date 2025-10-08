package handlers

var (
	DiscordBotConfigValues BotConfig
)

type AuthConfig struct {
	Token string `yaml:"token"`
}

type ChannelsConfig struct {
	PictureChannel string `yaml:"picture-channel"`
}

type StarConfig struct {
	Threshold int `yaml:"threshold"`
	Done      int `yaml:"done"`
}

type SettingsConfig struct {
	Star StarConfig `yaml:"star"`
}

type DiscordConfig struct {
	GuildIDs []string       `yaml:"guildIDs"`
	Auth     AuthConfig     `yaml:"auth"`
	Channels ChannelsConfig `yaml:"channels"`
	Redis    RedisConfig    `yaml:"redis"`
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

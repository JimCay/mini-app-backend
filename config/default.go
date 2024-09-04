package config

var defaultConfig = Config{
	TgConfig{},
	DbConfig{},
	LogConfig{
		Level: 4,
		Size:  "30m",
		Path:  "logs/",
	},
}

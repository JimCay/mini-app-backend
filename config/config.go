package config

import (
	"github.com/go-ini/ini"
	"github.com/urfave/cli/v2"
)

const DefaultConfigPath = "./config.ini"

type Config struct {
	TgConf  TgConfig
	DbConf  DbConfig
	LogConf LogConfig
}

type TgConfig struct {
	TelegramBotToken   string
	TelegramMiniAppUrl string
	Test               bool
	Ssl                bool
	Port               string
	Swagger            bool
}

type LogConfig struct {
	Level uint32
	Size  string
	Path  string
}

type DbConfig struct {
	Type        string
	User        string
	Password    string
	IP          string
	Port        string
	Name        string
	NumberShard int
}

func Setup(ctx *cli.Context) (*Config, error) {
	fig := &defaultConfig
	path := DefaultConfigPath
	if file := ctx.String(ConfigFileFlag.Name); file != "" {
		path = file
	}
	err := loadConfig(path, fig)
	if err != nil {
		return fig, err
	}
	return fig, nil
}

func loadConfig(filePath string, config *Config) error {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return err
	}

	mapTo(cfg, "telegram", &config.TgConf)
	mapTo(cfg, "db", &config.DbConf)
	mapTo(cfg, "log", &config.LogConf)

	return nil
}

func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		println("Cfg.MapTo %s err: %v", section, err)
	}
}

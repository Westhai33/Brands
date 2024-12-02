package config

import logger "Brands/pkg/zerohook"

type Config struct {
	Log      logger.LoggerConfig `yaml:"log"`
	Postgres struct {
		Conn string `yaml:"conn"`
	} `yaml:"postgres"`
}

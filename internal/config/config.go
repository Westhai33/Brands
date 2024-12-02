package config

import logger "Brands/pkg/zerohook"

type Config struct {
	Log      logger.LoggerConfig `yaml:"log"`
	Postgres struct {
		Conn string `yaml:"conn"`
	} `yaml:"postgres"`
	GRPCServer struct {
		Address string `yaml:"address"` // Адрес для gRPC сервера
	} `yaml:"grpcServer"`
	HTTPServer struct {
		Address string `yaml:"address"` // Адрес для HTTP сервера (Gateway)
	} `yaml:"httpServer"`
}

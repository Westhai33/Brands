package config

import logger "Brands/pkg/zerohook"

type Config struct {
	Log      logger.LoggerConfig `yaml:"log"`
	Postgres struct {
		Conn string `yaml:"conn"`
	} `yaml:"postgres"`
	Jaeger struct {
		AgentHost   string `yaml:"agent_host"`
		AgentPort   int    `yaml:"agent_port"`
		ServiceName string `yaml:"service_name"`
	} `yaml:"jaeger"`
	Prometheus struct {
		Port           int    `yaml:"port"`
		MetricsPath    string `yaml:"metrics_path"`
		ScrapeInterval string `yaml:"scrape_interval"`
	} `yaml:"prometheus"`
}

var (
	CorsAllowHeaders = "Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Max-Age, Access-Control-Allow-Credentials, Content-Type, Authorization, Origin, X-Requested-With , Accept"
	CorsAllowMethods = "HEAD, GET, POST, PUT, DELETE, OPTIONS"
	CorsAllowOrigin  = "*"
)

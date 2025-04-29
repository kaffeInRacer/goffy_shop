package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Env struct {
		Mode       string `yaml:"mode"`
		TimeZone   string `yaml:"time_zone"`
		TimeFormat string `yaml:"time_format"`
	}
	Server struct {
		HTTP struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"http"`
		GRPC struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"grpc"`
	} `yaml:"server"`
	Database struct {
		Postgres struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			DB       string `yaml:"db"`
		} `yaml:"postgres"`
		Redis struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		} `yaml:"redis"`
	} `yaml:"database"`
	Logger struct {
		LogLevel    string `yaml:"log_level"`
		LogFileName string `yaml:"log_file_name"`
	} `yaml:"logger"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

var Conf *config

func NewConfig(path string) {
	file, _ := os.Open(path)

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	_ = decoder.Decode(&Conf)

	if Conf.Env.Mode == "local" {
		// default setup logger
		Conf.Logger.LogLevel = "debug"
		Conf.Logger.LogFileName = "app_logger_local.log"

		// default HTTP
		Conf.Server.HTTP.Port = 8080
		Conf.Server.HTTP.Host = "localhost"

		// default gRPC
		Conf.Server.GRPC.Host = "localhost"
		Conf.Server.GRPC.Port = 9090

		// default postgres
		Conf.Database.Postgres.Host = "localhost"
		Conf.Database.Postgres.Port = 5432
		Conf.Database.Postgres.DB = "postgres"
		Conf.Database.Postgres.User = "postgres"
		Conf.Database.Postgres.Password = "postgres"

		// default redis
		Conf.Database.Redis.Host = "localhost"
		Conf.Database.Redis.Port = 6372
		Conf.Database.Redis.Password = "redis"
		Conf.Database.Redis.DB = 0
	}
}

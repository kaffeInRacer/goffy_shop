package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env struct {
		Mode       string `yaml:"mode"`
		TimeZone   string `yaml:"time_zone"`
		TimeFormat string `yaml:"time_format"`
	} `yaml:"env"`

	Logger struct {
		LogLevel    string `yaml:"log_level"`
		LogFileName string `yaml:"log_file_name"`
	} `yaml:"logger"`

	Server struct {
		HTTP struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"net_http"`
		GRPC struct {
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"gRPC"`
	} `yaml:"server"`

	Security struct {
		JWT struct {
			Secret string `yaml:"secret"`
		} `yaml:"jwt"`
	} `yaml:"security"`

	Database struct {
		Postgres struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Name     string `yaml:"name"`
			SSL      string `yaml:"ssl"`
		} `yaml:"postgres"`
		Redis struct {
			Host         string `yaml:"host"`
			Port         int    `yaml:"port"`
			Password     string `yaml:"password"`
			DB           int    `yaml:"db"`
			ReadTimeout  string `yaml:"read_timeout"`
			WriteTimeout string `yaml:"write_timeout"`
		} `yaml:"redis"`
	} `yaml:"database"`
}

var Conf *Config

func LoadConfig(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return err
	}

	Conf = &cfg
	return nil
}

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Address string `yaml:"address"`
	ApiKey  string `yaml:"api_key"`
}

func New() *Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

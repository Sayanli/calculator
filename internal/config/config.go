package config

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log  `yaml:"log"`
		HTTP `yaml:"http"`
		GRPC `yaml:"grpc"`
	}
	Log struct {
		Env string `yaml:"env"`
	}
	HTTP struct {
		Port int `yaml:"port"`
	}
	GRPC struct {
		Port int `yaml:"port"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join(configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}

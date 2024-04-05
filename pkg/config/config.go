package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env             string `yaml:"env" envDefault:"local"`
	DbFile          string `yaml:"db_file" env-required:"true"`
	SourceUrl       string `yaml:"source_url" envDefault:"https://xkcd.com"`
	SpellcheckModel string `yaml:"spellcheck_model" env:"SPELLCHECK_PATH" env-required:"true"`
}

func LoadConfig() (*Config, error) {

	// config path from env or default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &config, nil
}

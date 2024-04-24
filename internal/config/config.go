package config

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/internal/flags"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env       string `yaml:"env" envDefault:"local"`
	DBFile    string `yaml:"db_file" envDefault:"database.json"`
	SourceURL string `yaml:"source_url" envDefault:"https://xkcd.com"`
	Parallel  int    `yaml:"parallel" envRequired:"true"`
	BatchSize int    `yaml:"batch_size" envRequired:"true"`
}

func LoadConfig(flags *flags.UsersFlags) (*Config, error) {

	// config path from env, flag or try default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = flags.ConfigPath
	}
	if configPath == "" {
		configPath = "configs/config.yaml"
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

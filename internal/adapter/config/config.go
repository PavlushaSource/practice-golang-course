package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Args struct {
	configPath string
	port       string
	host       string
}

type Config struct {
	App struct {
		Env string `yaml:"env" envDefault:"local"`
	} `yaml:"App"`
	HTTP struct {
		Port           string        `yaml:"port" env-description:"Server port" env-default:"8080"`
		Host           string        `yaml:"host" env-description:"Server host" env-default:"localhost"`
		UpdateInterval time.Duration `yaml:"update-interval" envDefault:"24h"`
	} `yaml:"HTTP"`

	DB struct {
		Host        string `yaml:"host" env-description:"Database host"`
		Port        string `yaml:"port" env-description:"Database port"`
		Username    string `yaml:"username" env-description:"Database user name"`
		Password    string `env:"DB_PASSWORD" env-description:"Database user password"`
		Name        string `yaml:"db-name" env-description:"Database name"`
		Connections int    `yaml:"connections" env-description:"Total number of database connections"`
	} `yaml:"DB"`

	JSONFlat struct {
		DBFilepath    string `yaml:"file-path" env-description:"Path to comixs DB json file" envDefault:"database.json"`
		IndexFilepath string `yaml:"index-path" env-description:"Path to index comixs json file" envDefault:"index.json"`
	} `yaml:"JSONFlat"`

	ComixSource struct {
		URL       string `yaml:"url" env-description:"Comix source url" envDefault:"https://xkcd.com"`
		Parallel  int    `yaml:"parallel" envRequired:"true"`
		BatchSize int    `yaml:"batch_size" envRequired:"true"`
	} `yaml:"ComixSource"`
}

func New() (*Config, error) {

	args := getArgs()

	// configs path from env, flag or try default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = args.configPath
	}
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configs file not found: %w", err)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, fmt.Errorf("error reading configs file: %w", err)
	}

	if args.host != "" {
		config.HTTP.Host = args.host
	}
	if args.port != "" {
		config.HTTP.Port = args.port
	}

	return &config, nil
}

func getArgs() *Args {
	var configPath string
	var port string
	var host string
	flag.StringVar(&configPath, "c", "", "Path to configs.yaml")
	flag.StringVar(&port, "p", "", "Server port")
	flag.StringVar(&host, "h", "", "Server host")
	flag.Parse()
	return &Args{configPath, port, host}
}

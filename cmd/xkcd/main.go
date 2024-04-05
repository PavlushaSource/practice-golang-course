package main

import (
	"fmt"
	"github.com/PavlushaSource/yadro-practice-course/pkg/config"
	"os"
)

func main() {
	// read config file
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(cfg.Env)
}

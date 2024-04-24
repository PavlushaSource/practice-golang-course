package flags

import (
	"flag"
)

type UsersFlags struct {
	ConfigPath string
}

func GetFlagsFromCommandlineInput() *UsersFlags {
	var configPath string
	flag.StringVar(&configPath, "c", "", "Path to config.yaml")
	flag.Parse()
	return &UsersFlags{ConfigPath: configPath}
}

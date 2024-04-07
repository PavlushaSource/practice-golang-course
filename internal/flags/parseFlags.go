package flags

import (
	"flag"
	"fmt"
)

type UsersFlags struct {
	NumberComics    int
	OutputToConsole bool
}

func GetFlagsFromCommandlineInput() (*UsersFlags, error) {
	var numberComics int
	var outputToConsole bool
	flag.IntVar(&numberComics, "n", 0, "Enter an integer number for a specific number of comic book downloads")
	flag.BoolVar(&outputToConsole, "o", false, "Output comics saved in the database to the console")
	flag.Parse()
	if otherInput := flag.Args(); len(otherInput) > 0 {
		return nil, fmt.Errorf("xkcd app work works with specific arguments: -n and -o.\nEnter make help for more information")
	}
	return &UsersFlags{NumberComics: numberComics, OutputToConsole: outputToConsole}, nil
}

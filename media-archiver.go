package main

import (
	"fmt"
	"os"

	ma "github.com/asc8277/media-archiver/lib"
)

var Version = "dev"

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help())
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Println(Version)
	case "process":
		if len(os.Args) < 3 {
			fmt.Println(help())
			os.Exit(1)
		}
		archiver := ma.Archiver{InPath: os.Args[2]}
		archiver.Process()
	case "help":
		fmt.Println(help())
	default:
		fmt.Println(help())
		os.Exit(1)
	}
}

func help() string {
	return fmt.Sprintf("%s, valid commands: process <dir>, version, help", version())
}

func version() string {
	return fmt.Sprintf("media-archiver %s", Version)
}

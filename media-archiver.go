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

	switch os.Args[1] {
	case "--version":
		fmt.Println(version())
	case "--help":
		fmt.Println(help())
	default:
		archiver := ma.Archiver{InPath: os.Args[1]}
		archiver.Process()
	}
}

func help() string {
	return fmt.Sprintf("%s, usage: media-archiver <dir>, --version, --help", version())
}

func version() string {
	return fmt.Sprintf("media-archiver %s", Version)
}

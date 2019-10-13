package main

import (
	"os"

	ma "github.com/asc8277/media-archiver/lib"
)

func main() {
	inPath := os.Args[1]

	archiver := ma.Archiver{InPath: inPath}
	archiver.Process()
}

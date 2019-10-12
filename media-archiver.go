package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	mediaarchiver "github.com/asc8277/media-archiver/lib"
)

var supportedImages = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
}

var supportedVideos = map[string]bool{
	"mp4": true,
	"m4v": true,
}

func main() {
	inPath := os.Args[1]
	log.Print(inPath)

	files, err := ioutil.ReadDir(inPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var skipped []string
	var images []string
	var videos []string

	for _, file := range files {
		fname := file.Name()
		ext := mediaarchiver.File{Filename: fname, Path: inPath}.GetFileExtension()

		if supportedImages[ext] {
			images = append(images, fname)
		} else if supportedVideos[ext] {
			videos = append(videos, fname)
		} else {
			skipped = append(skipped, fname)
		}
	}

	outPath := filepath.Join(inPath, time.Now().Format("20060102_150405"))
	if len(images) > 0 || len(videos) > 0 {
		err = os.Mkdir(outPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	mediaarchiver.Files{Filenames: skipped, InPath: inPath, OutPath: outPath}.ProcessSkipped()
	log.Print("---")
	mediaarchiver.Files{Filenames: images, InPath: inPath, OutPath: outPath}.ProcessImages()
	log.Print("---")
	mediaarchiver.Files{Filenames: videos, InPath: inPath, OutPath: outPath}.ProcessVideos()
}

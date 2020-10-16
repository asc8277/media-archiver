package mediaarchiver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type MediaFormat int

const (
	Image MediaFormat = iota + 1
	Video
)

var supportedFormats = map[string]MediaFormat{
	"jpg":  Image,
	"jpeg": Image,
	"png":  Image,
	"mp4":  Video,
	"m4v":  Video,
	"mov": 	Video,
	"avi":	Video,
}

// Archiver archiver
type Archiver struct {
	InPath string
}

// Process process all files
func (ma *Archiver) Process() {
	files, err := ioutil.ReadDir(ma.InPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var skipped []string
	var images []string
	var videos []string

	for _, fileInfo := range files {
		fname := fileInfo.Name()
		mfile := file{name: fname, path: ma.InPath}
		ext := mfile.getFileExtension()

		if supportedFormats[ext] == Image {
			images = append(images, fname)
		} else if supportedFormats[ext] == Video {
			videos = append(videos, fname)
		} else {
			skipped = append(skipped, fname)
		}
	}

	outPath := filepath.Join(ma.InPath, time.Now().Format("20060102_150405"))
	if len(images) > 0 || len(videos) > 0 {
		err = os.Mkdir(outPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	log.Print(fmt.Sprintf("%s -> %s", ma.InPath, outPath))

	ma.processSkipped(skipped)
	log.Print("---")
	ma.processImages(images, ma.InPath, outPath)
	log.Print("---")
	ma.processVideos(videos, ma.InPath, outPath)
}

func (ma *Archiver) processSkipped(filenames []string) {
	for _, filename := range filenames {
		log.Print(fmt.Sprintf("%s: skipped", filename))
	}
}

func (ma *Archiver) processImages(filenames []string, inPath string, outPath string) {
	cpus := runtime.NumCPU()
	tasks := make(chan string)

	var wg sync.WaitGroup

	for worker := 0; worker < cpus; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for filename := range tasks {
				log.Print(fmt.Sprintf("%s: process", filename))

				imf := imageFile{in: file{name: filename, path: inPath}, out: file{path: outPath}}
				result := imf.process()

				log.Print(fmt.Sprintf("%s -> %s : %s", filename, imf.out.name, result))
			}
		}()
	}

	for _, task := range filenames {
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

func (ma *Archiver) processVideos(filenames []string, inPath string, outPath string) {
	for _, filename := range filenames {
		log.Print(fmt.Sprintf("%s: process", filename))

		vmf := videoFile{in: file{name: filename, path: inPath}, out: file{path: outPath}}
		result := vmf.process()

		log.Print(fmt.Sprintf("%s -> %s : %s", filename, vmf.out.name, result))
	}
}

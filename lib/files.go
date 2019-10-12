package mediaarchiver

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Files files
type Files struct {
	Filenames []string
	InPath    string
	OutPath   string
}

// ProcessSkipped skipped
func (files Files) ProcessSkipped() {
	for _, filename := range files.Filenames {
		log.Print(fmt.Sprintf("%s: skipped", filename))
	}
}

// ProcessImages images
func (files Files) ProcessImages() {
	cpus := runtime.NumCPU()
	tasks := make(chan string)

	var wg sync.WaitGroup

	for worker := 0; worker < cpus; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for filename := range tasks {
				mediafile := MediaFile{In: File{Filename: filename, Path: files.InPath}, Out: File{Path: files.OutPath}}

				fInPath := filepath.Join(files.InPath, mediafile.In.Filename)
				mediafile.SetNewImageFilename()
				fOutPath := filepath.Join(mediafile.Out.Path, mediafile.Out.Filename)

				log.Print(fmt.Sprintf("%s: -> %s", filename, fOutPath))

				out, err := exec.Command("jpeg-recompress", "-a", "-q", "high", "-n", "60", "-x", "95", fInPath, fOutPath).CombinedOutput()

				if err != nil {
					log.Fatal(err)
				}

				result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")

				log.Print(fmt.Sprintf("%s: %s %s", filename, result[len(result)-3], result[len(result)-2]))
			}
		}()
	}

	for _, task := range files.Filenames {
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

// ProcessVideos videos
func (files Files) ProcessVideos() {
	for _, filename := range files.Filenames {
		mediafile := MediaFile{In: File{Filename: filename, Path: files.InPath}, Out: File{Path: files.OutPath}}

		fInPath := filepath.Join(mediafile.In.Path, mediafile.In.Filename)
		mediafile.SetNewVideoFilename()
		fOutPath := filepath.Join(mediafile.Out.Path, mediafile.Out.Filename)

		log.Print(fmt.Sprintf("%s: -> %s", filename, fOutPath))

		out, err := exec.Command("HandBrakeCLI", "-i", fInPath, "-o", fOutPath, "-e", "x265", "-q", "22", "-f", "av_mp4", "--comb-detect", "--decomb", "-a", "1", "-E", "copy:aac", "--loose-anamorphic").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")

		log.Print(fmt.Sprintf("%s: %s", filename, result[len(result)-9]))
	}
}

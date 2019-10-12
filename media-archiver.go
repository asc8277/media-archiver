package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
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

func getFilePrefixFromFilename(fname string) string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d_\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(fname)
}

func getPartFilePrefixFromFilename(fname string) string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(fname)
}

func getFilePrefixFromExif(path string) string {
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
		return ""
	}

	x, err := exif.Decode(f)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	tm, err := x.DateTime()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return tm.Format("20060102_150405")
}

func getNewImageFilename(path string, fname string) string {
	fext := filepath.Ext(fname)
	fpre := filepath.Base(fname[0 : len(fname)-len(fext)])

	prefix := getFilePrefixFromExif(filepath.Join(path, fname))
	if prefix == "" {
		prefix = getFilePrefixFromFilename(fname)
	}
	if prefix == "" {
		prefix = getPartFilePrefixFromFilename(fname)
	}

	return fmt.Sprintf("%s_%s.%s", prefix, fpre, "jpg")
}

func getNewVideoFilename(path string, fname string) string {
	fext := filepath.Ext(fname)
	fpre := filepath.Base(fname[0 : len(fname)-len(fext)])

	prefix := getFilePrefixFromFilename(fname)
	if prefix == "" {
		prefix = getPartFilePrefixFromFilename(fname)
	}

	return fmt.Sprintf("%s_%s.%s", prefix, fpre, "mp4")
}

func processImages(fnames []string, path string, dPath string) {
	cpus := runtime.NumCPU()
	tasks := make(chan string)

	var wg sync.WaitGroup

	for worker := 0; worker < cpus; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for fname := range tasks {
				fpath := filepath.Join(path, fname)
				dfname := getNewImageFilename(path, fname)
				dfpath := filepath.Join(dPath, dfname)
				log.Print(fmt.Sprintf("%s: -> %s", fname, dfpath))
				out, err := exec.Command("jpeg-recompress", "-a", "-q", "high", "-n", "60", "-x", "95", fpath, dfpath).CombinedOutput()
				if err != nil {
					log.Fatal(err)
				}

				result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
				log.Print(fmt.Sprintf("%s: %s %s", fname, result[len(result)-3], result[len(result)-2]))
			}
		}()
	}

	for _, task := range fnames {
		tasks <- task
	}
	close(tasks)
	wg.Wait()
}

func processVideos(fnames []string, path string, dPath string) {
	for _, fname := range fnames {
		fpath := filepath.Join(path, fname)
		dfname := getNewVideoFilename(path, fname)
		dfpath := filepath.Join(dPath, dfname)
		log.Print(fmt.Sprintf("%s: -> %s", fname, dfpath))
		out, err := exec.Command("HandBrakeCLI", "-i", fpath, "-o", dfpath, "-e", "x265", "-q", "22", "-f", "av_mp4", "--comb-detect", "--decomb", "-a", "1", "-E", "copy:aac", "--loose-anamorphic").CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
		log.Print(fmt.Sprintf("%s: %s", fname, result[len(result)-9]))
	}
}

func processSkipped(paths []string) {
	for _, path := range paths {
		log.Print(fmt.Sprintf("%s: skipped", path))
	}
}

func getFileExtension(fpath string) string {
	fpathExt := filepath.Ext(fpath)
	ext := ""
	if len(fpathExt) > 0 {
		ext = strings.ToLower(fpathExt[1:])
	}

	return ext
}

func main() {
	path := os.Args[1]
	log.Print(path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	var skipped []string
	var images []string
	var videos []string

	for _, file := range files {
		fname := file.Name()
		ext := getFileExtension(fname)

		if supportedImages[ext] {
			images = append(images, fname)
		} else if supportedVideos[ext] {
			videos = append(videos, fname)
		} else {
			skipped = append(skipped, fname)
		}
	}

	dPath := filepath.Join(path, time.Now().Format("20060102_150405"))
	if len(images) > 0 || len(videos) > 0 {
		err = os.Mkdir(dPath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	processSkipped(skipped)
	log.Print("---")
	processImages(images, path, dPath)
	log.Print("---")
	processVideos(videos, path, dPath)
}

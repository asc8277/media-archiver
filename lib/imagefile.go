package mediaarchiver

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

type imageFile mediaFile

func (imf *imageFile) process() string {
	if imf.in.getFileExtension() == "png" {
		imf.preProcessPng()
	}
	fInPath := imf.in.getFullPath()
	imf.setNewFilename()
	fOutPath := imf.out.getFullPath()

	out, err := exec.Command("jpeg-recompress", "-a", "-q", "high", "-n", "60", "-x", "95", fInPath, fOutPath).CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
	return fmt.Sprintf("%s %s", result[len(result)-3], result[len(result)-2])
}

func (imf *imageFile) setNewFilename() imageFile {
	fpre := imf.in.getFileNameWithoutExtension()

	prefix := imf.getFilePrefixFromExif()
	if prefix == "" {
		prefix = imf.in.getFilePrefixFromFilename()
	}
	if prefix == "" {
		prefix = imf.in.getPartFilePrefixFromFilename()
	}

	if strings.HasPrefix(fpre, prefix) {
		imf.out.name = fmt.Sprintf("%s.%s", fpre, "jpg")
	} else {
		imf.out.name = fmt.Sprintf("%s_%s.%s", prefix, fpre, "jpg")
	}

	return *imf
}

func (imf *imageFile) getFilePrefixFromExif() string {
	f, err := os.Open(imf.in.getFullPath())

	if err != nil {
		return ""
	}

	x, err := exif.Decode(f)
	if err != nil {
		return ""
	}

	tm, err := x.DateTime()
	if err != nil {
		return ""
	}

	return tm.Format("20060102_150405")
}

func (imf *imageFile) preProcessPng() imageFile {
	f, _ := os.Open(imf.in.getFullPath())
	t, _, _ := image.Decode(f)

	filename := imf.in.getFileNameWithoutExtension() + ".jpg"
	f, _ = os.Create(filepath.Join(imf.in.path, filename))

	jpeg.Encode(f, t, &jpeg.Options{Quality: 100})

	imf.in.name = filename

	return *imf
}

package mediaarchiver

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/rwcarlsen/goexif/exif"
)

// MediaFile media file
type MediaFile struct {
	In  File
	Out File
}

// SetNewVideoFilename new video filename
func (file MediaFile) SetNewVideoFilename() MediaFile {
	fext := filepath.Ext(file.In.Filename)
	fpre := filepath.Base(file.In.Filename[0 : len(file.In.Filename)-len(fext)])

	prefix := file.getFilePrefixFromFilename()
	if prefix == "" {
		prefix = file.getPartFilePrefixFromFilename()
	}

	file.Out.Filename = fmt.Sprintf("%s_%s.%s", prefix, fpre, "mp4")

	return file
}

// SetNewImageFilename new image filename
func (file MediaFile) SetNewImageFilename() MediaFile {
	fext := filepath.Ext(file.In.Filename)
	fpre := filepath.Base(file.In.Filename[0 : len(file.In.Filename)-len(fext)])

	prefix := file.getFilePrefixFromExif()
	if prefix == "" {
		prefix = file.getFilePrefixFromFilename()
	}
	if prefix == "" {
		prefix = file.getPartFilePrefixFromFilename()
	}

	file.Out.Filename = fmt.Sprintf("%s_%s.%s", prefix, fpre, "jpg")

	return file
}

func (file MediaFile) getFilePrefixFromFilename() string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d_\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(file.In.Filename)
}

func (file MediaFile) getPartFilePrefixFromFilename() string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(file.In.Filename)
}

func (file MediaFile) getFilePrefixFromExif() string {
	f, err := os.Open(filepath.Join(file.In.Path, file.In.Filename))

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

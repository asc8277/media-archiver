package mediaarchiver

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

// mediaFile media file
type mediaFile struct {
	in  file
	out file
}

func (mf *mediaFile) processVideo() string {
	fInPath := filepath.Join(mf.in.path, mf.in.name)
	mf.setNewVideoFilename()
	fOutPath := filepath.Join(mf.out.path, mf.out.name)

	out, err := exec.Command("HandBrakeCLI", "-i", fInPath, "-o", fOutPath, "-e", "x265", "-q", "22", "-f", "av_mp4", "--comb-detect", "--decomb", "-a", "1", "-E", "copy:aac", "--loose-anamorphic").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
	return result[len(result)-9]
}

func (mf *mediaFile) processImage() string {
	fInPath := filepath.Join(mf.in.path, mf.in.name)
	mf.setNewImageFilename()
	fOutPath := filepath.Join(mf.out.path, mf.out.name)

	out, err := exec.Command("jpeg-recompress", "-a", "-q", "high", "-n", "60", "-x", "95", fInPath, fOutPath).CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}

	result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
	return fmt.Sprintf("%s %s", result[len(result)-3], result[len(result)-2])
}

// SetNewVideoFilename new video filename
func (mf *mediaFile) setNewVideoFilename() mediaFile {
	fext := filepath.Ext(mf.in.name)
	fpre := filepath.Base(mf.in.name[0 : len(mf.in.name)-len(fext)])

	prefix := mf.getFilePrefixFromFilename()
	if prefix == "" {
		prefix = mf.getPartFilePrefixFromFilename()
	}

	mf.out.name = fmt.Sprintf("%s_%s.%s", prefix, fpre, "mp4")

	return *mf
}

// SetNewImageFilename new image filename
func (mf *mediaFile) setNewImageFilename() mediaFile {
	fext := filepath.Ext(mf.in.name)
	fpre := filepath.Base(mf.in.name[0 : len(mf.in.name)-len(fext)])

	prefix := mf.getFilePrefixFromExif()
	if prefix == "" {
		prefix = mf.getFilePrefixFromFilename()
	}
	if prefix == "" {
		prefix = mf.getPartFilePrefixFromFilename()
	}

	mf.out.name = fmt.Sprintf("%s_%s.%s", prefix, fpre, "jpg")

	return *mf
}

func (mf *mediaFile) getFilePrefixFromFilename() string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d_\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(mf.in.name)
}

func (mf *mediaFile) getPartFilePrefixFromFilename() string {
	r, err := regexp.Compile(`\d\d\d\d\d\d\d\d`)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return r.FindString(mf.in.name)
}

func (mf *mediaFile) getFilePrefixFromExif() string {
	f, err := os.Open(filepath.Join(mf.in.path, mf.in.name))

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

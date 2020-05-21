package mediaarchiver

import (
	"path/filepath"
	"regexp"
	"strings"
)

type file struct {
	name string
	path string
}

type mediaFile struct {
	in  file
	out file
}

func (file *file) getFileExtension() string {
	fpathExt := filepath.Ext(file.name)

	ext := ""
	if len(fpathExt) > 0 {
		ext = strings.ToLower(fpathExt[1:])
	}

	return ext
}

func (file *file) getFullPath() string {
	return filepath.Join(file.path, file.name)
}

func (file *file) getFileNameWithoutExtension() string {
	return filepath.Base(file.name[0 : len(file.name)-len(file.getFileExtension())])
}

func (file *file) getFilePrefixFromFilename() string {
	r, _ := regexp.Compile(`\d\d\d\d\d\d\d\d_\d\d\d\d\d\d`)

	return r.FindString(file.name)
}

func (file *file) getPartFilePrefixFromFilename() string {
	r, _ := regexp.Compile(`\d\d\d\d\d\d\d\d`)

	return r.FindString(file.name)
}

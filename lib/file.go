package mediaarchiver

import (
	"path/filepath"
	"strings"
)

type file struct {
	name string
	path string
}

func (file *file) getFileExtension() string {
	fpathExt := filepath.Ext(file.name)

	ext := ""
	if len(fpathExt) > 0 {
		ext = strings.ToLower(fpathExt[1:])
	}

	return ext
}

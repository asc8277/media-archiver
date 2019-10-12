package mediaarchiver

import (
	"path/filepath"
	"strings"
)

// File file
type File struct {
	Filename string
	Path     string
}

// GetFileExtension get file extension
func (file File) GetFileExtension() string {
	fpathExt := filepath.Ext(file.Filename)

	ext := ""
	if len(fpathExt) > 0 {
		ext = strings.ToLower(fpathExt[1:])
	}

	return ext
}

package mediaarchiver

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type videoFile mediaFile

func (vmf *videoFile) process() string {
	fInPath := vmf.in.getFullPath()
	vmf.setNewFilename()
	fOutPath := vmf.out.getFullPath()

	out, err := exec.Command("HandBrakeCLI", "-i", fInPath, "-o", fOutPath, "-e", "x264", "-q", "23", "-f", "av_mp4", "--comb-detect", "--decomb", "-a", "1", "-E", "copy:aac", "--loose-anamorphic").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	result := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
	return result[len(result)-9]
}

// SetNewVideoFilename new video filename
func (vmf *videoFile) setNewFilename() videoFile {
	fpre := vmf.in.getFileNameWithoutExtension()

	prefix := vmf.in.getFilePrefixFromFilename()
	if prefix == "" {
		prefix = vmf.in.getPartFilePrefixFromFilename()
	}

	if strings.HasPrefix(fpre, prefix) {
		vmf.out.name = fmt.Sprintf("%s.%s", fpre, "mp4")
	} else {
		vmf.out.name = fmt.Sprintf("%s_%s.%s", prefix, fpre, "mp4")
	}

	return *vmf
}

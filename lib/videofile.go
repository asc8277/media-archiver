package mediaarchiver

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	cp "github.com/nmrshll/go-cp"
)

type videoFile mediaFile

func (vmf *videoFile) process() string {
	fInPath := vmf.in.getFullPath()
	vmf.setNewFilename()
	fOutPath := vmf.out.getFullPath()

	out, err := exec.Command("HandBrakeCLI", "-i", fInPath, "-o", fOutPath, "-e", "x264", "-q", "23", "-f", "av_mp4", "--comb-detect", "--decomb", "-a", "1", "-E", "copy:aac", "--loose-anamorphic").CombinedOutput()
	if err != nil {
		return err.Error()
	}

	resultArr := strings.Split(strings.ReplaceAll(string(out), "\r", ""), "\n")
	result := resultArr[len(resultArr)-9]

	if vmf.in.getFileExtension() == "mp4" {
		result = result + vmf.checkFileSizes()
	}

	return result
}

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

func (vmf *videoFile) checkFileSizes() string {
	fInPath := vmf.in.getFullPath()
	fOutPath := vmf.out.getFullPath()

	inFi, _ := os.Stat(fInPath)
	inSize := inFi.Size()

	outFi, _ := os.Stat(fOutPath)
	outSize := outFi.Size()

	if outSize > inSize {
		os.Remove(fOutPath)
		err := cp.CopyFile(fInPath, fOutPath)

		return err.Error()
	}

	return ""
}

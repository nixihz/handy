package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRunPath() string {
	workDir := os.Getenv("RICH_WORK_DIR")
	if workDir == "" {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		index := strings.LastIndex(path, string(os.PathSeparator))
		return path[:index]
	}

	return workDir
}

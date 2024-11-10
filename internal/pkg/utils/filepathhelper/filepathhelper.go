package filepathhelper

import (
	"os"
	"path/filepath"
)

func GetWorkDir() string {
	wDirPath, err := os.Getwd()
	if err != nil {
		if ex, er := os.Executable(); er == nil {
			wDirPath = filepath.Dir(ex)
		}
	}

	return wDirPath
}

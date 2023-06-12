package exe

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Path returns the current executable path
func Path() string {
	exePath := getPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(exePath, tmpDir) {
		return getPathByCaller()
	}
	return exePath
}

// getPathByExecutable returns the path by executable
func getPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// getPathByCaller returns the path by caller
func getPathByCaller() string {
	var exePath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		exePath = path.Dir(filename)
	}

	return exePath
}

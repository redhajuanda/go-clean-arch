package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// GetSourcePath returns source path
func GetSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func Dir() {
	// from Executable Directory
	ex, _ := os.Executable()
	fmt.Println("Executable DIR:", filepath.Dir(ex))

	// Current working directory
	dir, _ := os.Getwd()
	fmt.Println("CWD:", dir)

	// Relative on runtime DIR:
	_, b, _, _ := runtime.Caller(0)
	d1 := path.Join(path.Dir(b))
	fmt.Println("Relative", d1)
}

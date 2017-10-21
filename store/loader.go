package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const DirCount = 2

type Loader struct {
	file         string
	baseDir      string
	storeDirs    []string
	currentIndex int
	logger       *log.Logger
}

func NewLoader(file string, logger *log.Logger) *Loader {
	baseDir := filepath.Dir(file)
	return &Loader{
		file,
		baseDir,
		nil,
		0,
		logger,
	}
}

func (l *Loader) BuildStoreDirs() error {
	storeDirs, err := BuildDirs(l.baseDir, DirCount)
	if err != nil {
		return err
	}
	l.storeDirs = storeDirs
	return nil
}

func BuildDirs(baseDir string, count int) ([]string, error) {
	dirs := make([]string, count)

	for i := 0; i < DirCount; i++ {
		dir := filepath.Join(baseDir, fmt.Sprintf("db0%d", i))
		err := os.MkdirAll(dir, os.ModeDir)
		if err != nil {
			return nil, err
		}
		dirs[i] = dir
	}

	return dirs, nil
}

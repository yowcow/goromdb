package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const DirCount = 2
const DirPerm = 0755

type Loader struct {
	baseDir      string
	storeDirs    []string
	currentIndex int
	logger       *log.Logger
}

func NewLoader(baseDir string, logger *log.Logger) *Loader {
	storeDirs := make([]string, DirCount)
	return &Loader{
		baseDir,
		storeDirs,
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

func (l *Loader) MoveFileToNextDir(file string) (string, error) {
	nextIndex := l.currentIndex + 1
	if nextIndex == len(l.storeDirs) {
		nextIndex = 0
	}
	nextDir := l.storeDirs[nextIndex]
	base := filepath.Base(file)
	nextFile := filepath.Join(nextDir, base)
	if err := os.Rename(file, nextFile); err != nil {
		return "", err
	}
	l.currentIndex = nextIndex
	return nextFile, nil
}

func (l Loader) CleanOldDirs() error {
	for i, dir := range l.storeDirs {
		if i != l.currentIndex {
			if err := os.RemoveAll(dir); err != nil {
				return err
			}
			if err := os.MkdirAll(dir, DirPerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func BuildDirs(baseDir string, count int) ([]string, error) {
	dirs := make([]string, count)

	for i := 0; i < DirCount; i++ {
		dir := filepath.Join(baseDir, fmt.Sprintf("db0%d", i))
		if err := os.MkdirAll(dir, DirPerm); err != nil {
			return nil, err
		}
		dirs[i] = dir
	}

	return dirs, nil
}

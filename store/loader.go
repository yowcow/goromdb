package store

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// DirCount defines the number of subdirectories
const DirCount = 2

// DirPerm defines directory permission
const DirPerm = 0755

// Loader represents a loader
type Loader struct {
	basedir  string
	dirs     []string
	curindex int
}

// NewLoader creates a new loader
func NewLoader(basedir string) (*Loader, error) {
	dirs, err := buildDirs(basedir, DirCount)
	if err != nil {
		return nil, err
	}
	return &Loader{
		basedir,
		dirs,
		0,
	}, nil
}

func buildDirs(basedir string, count int) ([]string, error) {
	fi, err := os.Stat(basedir)
	if err != nil {
		return nil, err
	}
	if fi != nil && !fi.IsDir() {
		return nil, fmt.Errorf("path '%s' exists and not dir", basedir)
	}

	dirs := make([]string, count)
	for i := 0; i < count; i++ {
		dir := filepath.Join(basedir, fmt.Sprintf("data%02d", i))
		dirs[i] = dir
		if _, err := os.Stat(dir); err != nil {
			if err = os.Mkdir(dir, DirPerm); err != nil {
				return nil, err
			}
		}
	}
	return dirs, nil
}

// DropIn drops given file into next subdirectory, and returns the filepath
func (l *Loader) DropIn(file string) (string, error) {
	defer syscall.Sync() // make sure write is in sync

	nextindex := l.curindex + 1
	if nextindex >= len(l.dirs) {
		nextindex = 0
	}
	base := filepath.Base(file)
	nextdir := l.dirs[nextindex]
	nextfile := filepath.Join(nextdir, base)
	if err := os.Rename(file, nextfile); err != nil {
		return nextfile, err
	}
	curdir := l.dirs[l.curindex]
	curfile := filepath.Join(curdir, base)
	if _, err := os.Stat(curfile); err == nil {
		os.Remove(curfile)
	}
	l.curindex = nextindex
	return nextfile, nil
}

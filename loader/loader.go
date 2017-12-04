package loader

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
	basedir   string
	filename  string
	dirs      []string
	curindex  int
	previndex int
}

// NewLoader creates a new loader
func New(basedir, filename string) (*Loader, error) {
	dirs, err := buildDirs(basedir, DirCount)
	if err != nil {
		return nil, err
	}
	return &Loader{
		basedir,
		filename,
		dirs,
		-1,
		-1,
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

// FindAny tries finding a file to load in any existing subdirectiries, and returns its filepath
func (l Loader) FindAny() (string, bool) {
	for i := 0; i < DirCount; i++ {
		file := filepath.Join(l.dirs[i], l.filename)
		if _, err := os.Stat(file); err == nil {
			return file, true
		}
	}
	return "", false
}

// DropIn drops given file into next subdirectory, and returns the filepath
func (l *Loader) DropIn(file string) (string, error) {
	defer syscall.Sync() // make sure write is in sync

	nextindex := l.curindex + 1
	if nextindex >= len(l.dirs) {
		nextindex = 0
	}
	nextdir := l.dirs[nextindex]
	nextfile := filepath.Join(nextdir, l.filename)
	if err := os.Rename(file, nextfile); err != nil {
		return nextfile, err
	}
	l.previndex = l.curindex
	l.curindex = nextindex
	return nextfile, nil
}

// CleanUp cleans previously loaded data file, and returns bool
func (l Loader) CleanUp() bool {
	if l.previndex < 0 {
		return false
	}
	prevdir := l.dirs[l.previndex]
	prevfile := filepath.Join(prevdir, l.filename)
	if err := os.Remove(prevfile); err != nil {
		return false
	}
	return true
}

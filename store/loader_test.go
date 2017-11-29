package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/testutil"
)

func TestBuildDirs(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	type Case struct {
		basedir      string
		expectError  bool
		expectedDirs []string
		subtest      string
	}
	cases := []Case{
		{
			"/tmp/hoge/fuga",
			true,
			nil,
			"non-existing basedir fails",
		},
		{
			dir,
			false,
			[]string{
				filepath.Join(dir, "data00"),
				filepath.Join(dir, "data01"),
			},
			"existing basedir succeeds",
		},
		{
			dir,
			false,
			[]string{
				filepath.Join(dir, "data00"),
				filepath.Join(dir, "data01"),
			},
			"re-creating dirs succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			dirs, err := buildDirs(c.basedir, 2)

			if c.expectError {
				assert.NotNil(t, err)
				assert.Nil(t, dirs)
			} else {
				assert.Nil(t, err)
				for i, dir := range dirs {
					assert.Equal(t, c.expectedDirs[i], dir)
				}
			}
		})
	}
}

func TestNewLoader(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	loader, err := NewLoader(dir)

	assert.Nil(t, err)
	assert.NotNil(t, loader)
}

func TestDropIn(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	type Case struct {
		expectedFilepath string
		expectedNotExist string
		subtest          string
	}
	cases := []Case{
		{
			filepath.Join(dir, "data01", "dropped-in"),
			filepath.Join(dir, "data00", "dropped-in"),
			"1st drop-in stores into data01",
		},
		{
			filepath.Join(dir, "data00", "dropped-in"),
			filepath.Join(dir, "data01", "dropped-in"),
			"2nd drop-in stores into data00",
		},
		{
			filepath.Join(dir, "data01", "dropped-in"),
			filepath.Join(dir, "data00", "dropped-in"),
			"3rd drop-in stores into data01",
		},
	}

	loader, _ := NewLoader(dir)

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			input := filepath.Join(dir, "dropped-in")
			testutil.CopyFile(input, "loader_test.go")
			actual, err := loader.DropIn(input)

			assert.Nil(t, err)
			assert.Equal(t, c.expectedFilepath, actual)

			_, err = os.Stat(actual)

			assert.Nil(t, err)

			_, err = os.Stat(c.expectedNotExist)

			assert.True(t, os.IsNotExist(err))
		})
	}
}

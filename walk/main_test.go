package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		testName string
		cfg      config
		expected string
	}{
		{testName: "NoFilter", cfg: config{root: "testdata", list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{testName: "FilterExtensionMatch", cfg: config{root: "testdata",
			ext: ".log", list: true}, expected: "testdata/dir.log\n"},
		{testName: "FilterExtensionSizeMatch", cfg: config{root: "testdata",
			ext: "log", minSize: 10, list: true}, expected: "testdata/dir.log\n"},
		{testName: "FilterExtensionSizeNoMatch", cfg: config{root: "testdata",
			ext: "log", minSize: 20, list: true}, expected: ""},
		{testName: "FilterExtensionNoMatch", cfg: config{root: "testdata",
			ext: ".gz", minSize: 0, list: true}, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			buffer := bytes.Buffer{}

			err := run(&buffer, tc.cfg)
			assert.Nil(t, err)

			result := buffer.String()

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		testName        string
		cfg             config
		filesToCreate   map[string]int
		filesRemaining  map[string]int
		numFilesDeleted int
		expected        string
	}{
		{
			testName: "DeleteExtensionNoMatch",
			cfg:      config{ext: ".log", del: true},
			filesToCreate: map[string]int{
				".gz": 10,
			},
			filesRemaining: map[string]int{
				".gz": 10,
			},
			numFilesDeleted: 0,
			expected:        "",
		},
		{
			testName: "DeleteExtensionMatch",
			cfg:      config{ext: ".log", del: true},
			filesToCreate: map[string]int{
				".log": 10,
			},
			filesRemaining:  map[string]int{},
			numFilesDeleted: 10,
			expected:        "",
		},
		{
			testName: "DeleteExtensionMixed",
			cfg:      config{ext: ".log", del: true},
			filesToCreate: map[string]int{
				".gz":  5,
				".log": 5,
			},
			filesRemaining: map[string]int{
				".gz": 5,
			},
			numFilesDeleted: 5,
			expected:        "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tempDir, cleanup := createTempDir(t, tc.filesToCreate)
			defer cleanup()

			tc.cfg.root = tempDir
			tc.cfg.wLog = &logBuffer
			err := run(&buffer, tc.cfg)
			assert.Nil(t, err)

			result := buffer.String()
			assert.Equal(t, tc.expected, result)

			// Verify the remaining files

			filesRemaining, err := os.ReadDir(tempDir)
			assert.Nil(t, err)

			tempMap := map[string]int{}
			for _, file := range filesRemaining {
				ext := filepath.Ext(file.Name())
				tempMap[ext]++
			}

			assert.Equal(t, tc.filesRemaining, tempMap)

			// Verify the number of log lines
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			assert.Equal(t, tc.numFilesDeleted, len(lines)-1)
		})
	}
}

func createTempDir(t *testing.T,
	files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := ioutil.TempDir("", "walktest")
	assert.Nil(t, err)

	for ext, numFiles := range files {
		// create 'numFiles' number of files with extension 'ext'
		for n := 0; n < numFiles; n++ {
			fname := fmt.Sprintf("file%d%s", n, ext)
			fpath := filepath.Join(tempDir, fname)
			err := os.WriteFile(fpath, []byte("dummy"), 0644)
			assert.Nil(t, err)
		}
	}

	dirname = tempDir
	cleanup = func() {
		os.RemoveAll(tempDir)
	}
	return
}

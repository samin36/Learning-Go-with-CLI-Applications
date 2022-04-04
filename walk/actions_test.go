package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		testName  string
		fileName  string
		ext       string
		minSize   uint64
		filterOut bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionMatchNoDot", "testdata/dir.log", "log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterDirectoryNoExt", "testdata", "", 0, true},
		{"FilterDirectoryExt", "testdata", ".ext", 0, true},
		{"FilterSizeMatch", "testdata/dir.log", "", 10, false},
		{"FilterSizeMatchExt", "testdata/dir.log", ".log", 10, false},
		{"FilterSizeNoMatch", "testdata/dir.log", "", 20, true},
		{"FilterSizeNoMatchExt", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			info, err := os.Stat(tc.fileName)
			assert.Nil(t, err)

			result := filterOut(info, config{
				ext:     tc.ext,
				minSize: tc.minSize,
			})

			assert.Equal(t, tc.filterOut, result)
		})
	}
}

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/golden_test1.html"
)

func TestParseContent(t *testing.T) {
	input := readFile(t, inputFile)

	result := parseContent(input)

	expected := readFile(t, goldenFile)

	assert.Equal(t, expected, result)
}

func TestRun(t *testing.T) {
	var mockStdout bytes.Buffer

	if err := run(inputFile, &mockStdout); err != nil {
		t.Fatal(err)
	}

	saveMsg := mockStdout.String()
	toIndex := strings.Index(saveMsg, "to")
	assert.NotEqual(t, -1, toIndex)
	dotHtmlIndex := strings.Index(saveMsg, ".html")
	outFname := strings.TrimSpace(saveMsg[toIndex+3 : dotHtmlIndex+5])
	defer cleanup(outFname)

	result := readFile(t, outFname)

	expected := readFile(t, goldenFile)

	assert.Equal(t, expected, result)
}

func cleanup(outFname string) {
	os.Remove(outFname)
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

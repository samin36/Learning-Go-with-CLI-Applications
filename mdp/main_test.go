package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "./testdata/test1.html"
	goldenFile = "./testdata/golden_test1.html"
)

func TestParseContent(t *testing.T) {
	input := readFile(t, inputFile)

	result := parseContent(input)

	expected := readFile(t, goldenFile)

	assert.Equal(t, expected, result)
}

func TestRun(t *testing.T) {
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}

	result := readFile(t, resultFile)

	expected := readFile(t, goldenFile)

	assert.Equal(t, expected, result)
	cleanup()
}

func cleanup() {
	os.Remove(resultFile)
}

func readFile(t *testing.T, filename string) []byte {
	t.Helper()

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

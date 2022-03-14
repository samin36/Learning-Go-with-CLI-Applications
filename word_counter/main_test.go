package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {

	testCases := []struct {
		words    string
		numWords int
	}{
		{words: "word1 word2 word3 word4\n", numWords: 4},
		{words: "word1 word2 word3 word4", numWords: 4},
		{words: "word1 \n   \n  word2  \n  word3", numWords: 3},
		{words: "word1 \n   \n    \n  word2\n", numWords: 2},
		{words: "\nword1\n", numWords: 1},
		{words: "\n", numWords: 0},
		{words: "", numWords: 0},
	}

	for _, testCase := range testCases {
		b := bytes.NewBufferString(testCase.words)
		got := count(b, false)
		assertCount(t, testCase.words, false, got, testCase.numWords)
	}
}

func TestCountLines(t *testing.T) {

	testCases := []struct {
		lines    string
		numLines int
	}{
		{lines: "word1 \n\n  word2  \n  word3", numLines: 4},
		{lines: "word1 \n   \n   \n  word2\n", numLines: 4},
		{lines: "\nword1", numLines: 2},
		{lines: "\nword1\n", numLines: 2},
		{lines: "word1 word2 word3 word4", numLines: 1},
		{lines: "word1 word2 word3 word4\n", numLines: 1},
		{lines: "\n", numLines: 1},
		{lines: "", numLines: 0},
	}

	for _, testCase := range testCases {
		b := bytes.NewBufferString(testCase.lines)
		got := count(b, true)
		assertCount(t, testCase.lines, true, got, testCase.numLines)
	}
}

func assertCount(t testing.TB, str string, countLines bool, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("expected %d for count(%q, %v), got %d instead\n", want, str, countLines, got)
	}
}

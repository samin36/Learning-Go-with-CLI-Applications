// Main entry point package for the word counter CLI
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// Counts the number of words in the data read from 'r'
func count(r io.Reader, countLines bool) (wc int) {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	for scanner.Scan() {
		wc++
	}

	return wc
}

func main() {
	countLines := flag.Bool("l", false, "Count lines")

	flag.Parse()

	fmt.Println(count(os.Stdin, *countLines))
}

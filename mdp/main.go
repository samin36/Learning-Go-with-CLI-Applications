package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// color objects
var (
	errColor = color.New(color.FgRed, color.Bold)
	okColor  = color.New(color.FgGreen)
)

// blackfriday package doesn't include the HTML header and footer
const (
	header = `<!DOCTYPE html>
	<html>
		<head>
			<meta http-equiv="content-type" content="text/html; charset=utf-8">
			<title>Markdown Preview Tool</title>
		</head>
		<body>
	`
	footer = `
		</body>
	</html>
	`
)

// define ConstError
type ConstError string

func (ce ConstError) Error() string {
	return string(ce)
}

// define errors
const (
	ErrFileNotSpecified = ConstError("Markdown file not specified")
)

func main() {
	file := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		fmt.Println()
		exit(ErrFileNotSpecified)
	}

	if err := run(*file); err != nil {
		exit(err)
	}
}

func run(filename string) error {
	// Read the data from the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(data)

	dotMdIndex := strings.LastIndex(filename, ".md")
	if dotMdIndex == -1 {
		dotMdIndex = len(filename)
	}
	outFname := fmt.Sprintf("%s.html", filename[:dotMdIndex])
	err = saveHTML(outFname, htmlData)

	if err == nil {
		okColor.Printf("Successfully saved HTML to %s\n", outFname)
	}

	return err
}

func parseContent(data []byte) []byte {
	// Parse the markdown file
	output := blackfriday.Run(data)

	// Generate a valid and safe HTML
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	// Write the header + body + footer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outFname string, htmlData []byte) error {
	// Write the bytes to the HTML file
	return os.WriteFile(outFname, htmlData, 0644)
}

func exit(err error) {
	errColor.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	ErrOsNotSupported   = ConstError("OS not supported")
)

// define messages
type Msg string

const (
	SaveSuccessMsg = "Successfully saved HTML to %s\n"
)

func main() {
	file := flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *file == "" {
		flag.Usage()
		fmt.Println()
		exit(ErrFileNotSpecified)
	}

	if err := run(*file, os.Stdout); err != nil {
		exit(err)
	}
}

func run(filename string, dest io.Writer) error {
	// Read the data from the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(data)

	// Create a temporary fiile and check for errors
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outFname := temp.Name()
	err = saveHTML(outFname, htmlData)

	if err == nil {
		okColor.Fprintf(dest, SaveSuccessMsg, outFname)
	}

	defer os.Remove(outFname)

	return preview(outFname)
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

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	// Define executable based on the OS
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = append(cParams, "/C", "start")
	case "darwin":
		cName = "open"
	default:
		return ErrOsNotSupported
	}

	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	return exec.Command(cPath, cParams...).Run()
}

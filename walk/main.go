package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	root := flag.String("root", ".", "Root directory to start")
	logFile := flag.String("log", "", `Log deletes to this file. By default,`+
		`it will be sent to STDOUT`)
	// Action options
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	minSize := flag.Uint64("minSize", 0, "Minimum file size")
	flag.Parse()

	cfg := config{
		root:    *root,
		list:    *list,
		del:     *del,
		ext:     *ext,
		minSize: *minSize,
	}
	//configure the options
	exit(cfg.configure(*logFile))
	// verify the options
	exit(cfg.verify())

	// run the program
	exit(run(os.Stdout, cfg))
}

func run(out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	err := filepath.Walk(cfg.root,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(info, cfg) {
				return nil
			}

			switch {
			case cfg.list:
				return listFile(path, out)
			case cfg.del:
				return deleteFile(path, delLogger)
			default:
				return nil
			}
		})

	return err
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

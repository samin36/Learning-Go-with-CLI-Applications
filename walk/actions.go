package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func filterOut(info fs.FileInfo, cfg config) bool {
	// if the extension filter doesn't start with a dot, then add one
	if len(cfg.ext) > 0 && cfg.ext[0] != '.' {
		cfg.ext = "." + cfg.ext
	}

	switch {
	case info.IsDir():
	case len(cfg.ext) > 0 && filepath.Ext(info.Name()) != cfg.ext:
	case info.Size() < int64(cfg.minSize):
	default:
		return false
	}

	return true
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}

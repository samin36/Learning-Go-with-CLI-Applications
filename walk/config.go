package main

import (
	"fmt"
	"io"
	"os"
)

// custom configuration related errors
type ConfigError string

func (ce ConfigError) Error() string {
	return string(ce)
}

func (ce ConfigError) Errorf(args ...interface{}) error {
	return fmt.Errorf(string(ce), args...)
}

const (
	ErrDirNotFound = ConfigError("%s: directory not found")
)

// all the configuration options
type config struct {
	// root directory to start searching from
	root string
	// extension
	ext string
	// min file size
	minSize uint64
	// list files
	list bool
	// delete fies
	del bool
	// log destination writer
	wLog io.Writer
}

func (c *config) configure(logFile string) error {
	// Configure the log destination
	var (
		f         = os.Stdout
		err error = nil
	)
	if logFile != "" {
		f, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}

		defer f.Close()
	}
	c.wLog = f

	return nil
}

func (c *config) verify() error {
	// verify the root directory exists
	if _, err := os.Stat(c.root); err != nil {
		if os.IsNotExist(err) {
			return ErrDirNotFound.Errorf(c.root)
		}
	}

	return nil
}

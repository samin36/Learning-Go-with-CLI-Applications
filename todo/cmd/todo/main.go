package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"todo"
)

// Default file name
var todoFilename = ".todo.json"

type ConstError string

func (ce ConstError) Error() string {
	return string(ce)
}

// Errors
const (
	ErrInvalidOption = ConstError("invalid option provided")
	ErrMissingValue  = ConstError("value missing")
)

func main() {
	// Parse the command line flags
	task := flag.String("task", "", "Task to be added in the ToDo list. Specify 'STDIN' to supply task name using stdin")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	l := &todo.List{}

	// Get the TODO_FILENAME environment variable
	if val := os.Getenv("TODO_FILENAME"); len(val) > 0 {
		todoFilename = val
	}

	// Read any existing items from the file
	exit(l.Get(todoFilename))

	// Parse the args
	switch {
	case *list:
		// List current todo items
		fmt.Println(l)
	case *complete > 0:
		exit(l.Complete(*complete))

		// Save the list
		exit(l.Save(todoFilename))
	case *task == "STDIN":
		taskName, err := getTask(os.Stdin)

		if err != nil {
			exit(err)
		}

		l.Add(taskName)

		// Save the list
		exit(l.Save(todoFilename))
	case len(*task) > 0:
		l.Add(*task)

		// Save the list
		exit(l.Save(todoFilename))
	default:
		exit(ErrInvalidOption)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getTask(r io.Reader) (string, error) {
	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if text := s.Text(); len(text) == 0 {
		return "", ErrMissingValue
	} else {
		return text, nil
	}
}

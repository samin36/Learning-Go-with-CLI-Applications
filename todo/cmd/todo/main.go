package main

import (
	"flag"
	"fmt"
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
)

func main() {
	// Parse the command line flags
	task := flag.String("task", "", "Task to be added in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	l := &todo.List{}

	// Get the TODO_FILENAME environment variable
	if val := os.Getenv("TODO_FILENAME"); len(val) > 0 {
		todoFilename = val
	}

	// Read any existing items from the file
	Exit(l.Get(todoFilename))

	// Parse the args
	switch {
	case *list:
		// List current todo items
		fmt.Println(l)
	case *complete > 0:
		Exit(l.Complete(*complete))

		// Save the list
		Exit(l.Save(todoFilename))
	case len(*task) > 0:
		l.Add(*task)

		// Save the list
		Exit(l.Save(todoFilename))
	default:
		Exit(ErrInvalidOption)
	}
}

func Exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

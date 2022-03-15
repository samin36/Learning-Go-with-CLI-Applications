package main

import (
	"fmt"
	"os"
	"strings"
	"todo"
)

// Temporarily hardcode the file name
const todoFilename = ".todo.json"

func main() {
	l := &todo.List{}

	// Read any existing items from the file
	if err := l.Get(todoFilename); err != nil {
		Exit(err)
	}

	// Parse the args
	switch len(os.Args) {
	case 1:
		// No user args
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	default:
		item := strings.Join(os.Args[1:], " ")
		l.Add(item)

		if err := l.Save(todoFilename); err != nil {
			Exit(err)
		}
	}
}

func Exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

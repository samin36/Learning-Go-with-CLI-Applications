// Contains the core implementation for the TODO application, separate from any
// CLI logic.
package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type ConstError string

func (ce ConstError) Error() string {
	return string(ce)
}

func (ce ConstError) Errorf(args ...interface{}) error {
	return fmt.Errorf(ce.Error(), args)
}

const (
	ErrItemNotFound         = ConstError("item #%d does not exist")
	ErrItemAlreadyCompleted = ConstError("item #%d has already been completed")
)

// Represents a todo item.
// 'item' is private and not accessible by other packages
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// Represents a list of todo items
type List []item

// Create a new todo item and add an item to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Mark a todo item as completed.
// 'itemNum' is NOT 0 based. Starts at 1 and goes up to length of list
func (l *List) Complete(itemNum int) error {
	list := *l

	if itemNum <= 0 || itemNum > len(list) {
		return ErrItemNotFound.Errorf(itemNum)
	}

	// Note the '&'. If it is omitted, 'todoItem' will contain a copy of
	// the item struct at list[itemNum-1]
	todoItem := &list[itemNum-1]

	if todoItem.Done {
		return ErrItemAlreadyCompleted.Errorf(itemNum)
	}

	todoItem.Done = true
	todoItem.CompletedAt = time.Now()

	return nil
}

// Delete a todo item from the list.
// 'itemNum' is NOT 0 based. Starts at 1 and goes up to length of list
func (l *List) Delete(itemNum int) error {
	list := *l

	if itemNum <= 0 || itemNum > len(list) {
		return ErrItemNotFound.Errorf(itemNum)
	}

	index := itemNum - 1

	*l = append(list[:index], list[index+1:]...)

	return nil
}

// Save method encodes the List as JSON and saves it
// using the provided filename
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Open the file specified by the provided filename, decode
// the JSON data, and parse it into a List
func (l *List) Get(filename string) error {
	js, err := os.ReadFile(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return json.Unmarshal(js, l)
}

//String prints out a formatted list by implementing the
// fmt.Stringer interface
func (l *List) String() (formatted string) {
	for itemNum, item := range *l {
		prefix := "[ ] "
		if item.Done {
			prefix = "[X] "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, itemNum+1, item.Task)
	}

	return formatted
}

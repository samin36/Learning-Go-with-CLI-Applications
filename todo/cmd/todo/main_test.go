package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"todo"

	"github.com/stretchr/testify/assert"
)

type Cmd struct {
	Args []string
	Err  error
	Out  string
}

var (
	binName  = "todo"
	filename = ".todo.json"
)

const (
	LIST     = "-list"
	ADD      = "-task"
	COMPLETE = "-complete"

	TASKNAME1 = "item1"
	TASKNAME2 = "item2"
)

func TestMain(m *testing.M) {

	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(filename)

	os.Exit(result)
}

func TestTodoCli(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	testCases := []struct {
		testName string
		cmds     []Cmd
	}{
		{
			testName: "Add;List;Complete;List",
			cmds: []Cmd{
				{Args: []string{ADD, TASKNAME1}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: format(1, TASKNAME1, false, true)},
				{Args: []string{COMPLETE, "1"}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: format(1, TASKNAME1, true, true)},
			},
		},
		{
			testName: "List",
			cmds: []Cmd{
				{Args: []string{LIST}, Err: nil, Out: "\n"},
			},
		},
		{
			testName: "Complete",
			cmds: []Cmd{
				{Args: []string{COMPLETE, "1"}, Err: fmt.Errorf("exit status 1"), Out: todo.ErrItemNotFound.Errorf(1).Error() + "\n"},
			},
		},
		{
			testName: "Add1;Complete1;List;Add2;List;Complete1",
			cmds: []Cmd{
				{Args: []string{ADD, TASKNAME1}, Err: nil, Out: ""},
				{Args: []string{COMPLETE, "1"}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: format(1, TASKNAME1, true, true)},
				{Args: []string{ADD, "STDIN", TASKNAME2}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: format(1, TASKNAME1, true, false) + "\n" + format(2, TASKNAME2, false, true)},
				{Args: []string{COMPLETE, "1"}, Err: fmt.Errorf("exit status 1"), Out: todo.ErrItemAlreadyCompleted.Errorf(1).Error() + "\n"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			cmds := testCase.cmds

			for i, cmd := range cmds {
				actual_cmd := exec.Command(cmdPath, cmd.Args...)

				if cmd.Args[0] == ADD && cmd.Args[1] == "STDIN" {
					actual_cmd = exec.Command(cmdPath, cmd.Args[:2]...)
					cmdStdin, err := actual_cmd.StdinPipe()
					assert.Nil(t, err)

					_, err = io.WriteString(cmdStdin, cmd.Args[2])
					assert.Nil(t, err)
					cmdStdin.Close()
				}

				out, err := actual_cmd.CombinedOutput()
				defer os.Remove(filename)

				if cmd.Err == nil {
					assert.Nil(t, err, "Err")
				} else {
					assert.EqualError(t, err, cmd.Err.Error(), "Err")
				}

				assert.Equal(t, cmd.Out, string(out), fmt.Sprintf("cmd #%d: %q", i+1, cmd))
			}

		})
	}
}

func format(itemNum int, TASKNAME1 string, completed, appendNewLines bool) string {
	prefix := "[ ]"
	if completed {
		prefix = "[X]"
	}

	newLines := "\n\n"
	if !appendNewLines {
		newLines = ""
	}

	return fmt.Sprintf("%s %d: %s%s", prefix, itemNum, TASKNAME1, newLines)
}

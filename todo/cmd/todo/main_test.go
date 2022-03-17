package main_test

import (
	"fmt"
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

	TASKNAME = "item1"
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
				{Args: []string{ADD, TASKNAME}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: TASKNAME + "\n"},
				{Args: []string{COMPLETE, "1"}, Err: nil, Out: ""},
				{Args: []string{LIST}, Err: nil, Out: ""},
			},
		},
		{
			testName: "List",
			cmds: []Cmd{
				{Args: []string{LIST}, Err: nil, Out: ""},
			},
		},
		{
			testName: "Complete",
			cmds: []Cmd{
				{Args: []string{COMPLETE, "1"}, Err: fmt.Errorf("exit status 1"), Out: todo.ErrItemNotFound.Errorf(1).Error() + "\n"},
			},
		},
		{
			testName: "Add;Complete;Complete",
			cmds: []Cmd{
				{Args: []string{ADD, TASKNAME}, Err: nil, Out: ""},
				{Args: []string{COMPLETE, "1"}, Err: nil, Out: ""},
				{Args: []string{COMPLETE, "1"}, Err: fmt.Errorf("exit status 1"), Out: todo.ErrItemAlreadyCompleted.Errorf(1).Error() + "\n"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			cmds := testCase.cmds

			for i, cmd := range cmds {
				actual_cmd := exec.Command(cmdPath, cmd.Args...)

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

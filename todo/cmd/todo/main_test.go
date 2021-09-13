package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binaryName = "todo"
	fileName   = ".todo.json"
)

func TestMain(t *testing.M) {
	fmt.Println("building tool . . . ")

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binaryName)
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("running tests . . . ")

	result := t.Run()

	fmt.Println("cleaning up . . . ")
	os.Remove(binaryName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "one two three"

	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	cmdPath := filepath.Join(dir, binaryName)

	t.Run("add new task", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)

		if err := cmd.Run(); err != nil {
			t.Fatalf("could not run cmd: %s", err)
		}
	})

	task2 := "test task number 2"
	t.Run("add new task from stdin", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()

		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("list tasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("could not get output: %s", err)
		}

		want := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)

		if want != string(out) {
			t.Errorf("want %s but got %s", want, string(out))
		}
	})
}

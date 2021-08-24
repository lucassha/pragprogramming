package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
		cmd := exec.Command(cmdPath, strings.Split(task, " ")...)

		if err := cmd.Run(); err != nil {
			t.Fatalf("could not run cmd: %s", err)
		}
	})

	t.Run("list tasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath)

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("could not get output: %s", err)
		}

		want := task + "\n"

		if want != string(out) {
			t.Errorf("want %s but got %s", want, string(out))
		}
	})
}

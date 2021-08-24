package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lucassha/pragprogramming/todo"
)

const (
	todoFileName = ".todo.json"
)

func main() {
	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	// 1 arg means it's only the binary name passed in
	case len(os.Args) == 1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	default:
		task := strings.Join(os.Args[1:], " ")
		l.Add(task)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

}

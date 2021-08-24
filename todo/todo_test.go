package todo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lucassha/pragprogramming/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskName := "new task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("got %v but want %v", l[0].Task, taskName)
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}

	taskName := "new task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("got %v but want %v", l[0].Task, taskName)
	}

	if l[0].Done {
		t.Errorf("new task should not be completed")
	}

	l.Complete(1)

	if !l[0].Done {
		t.Errorf("task should be completed")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}

	tasks := []string{"task1", "task2", "task3"}

	for _, v := range tasks {
		l.Add(v)
	}

	if l[0].Task != tasks[0] {
		t.Errorf("want %q, got %q", tasks[0], l[0].Task)
	}

	l.Delete(2)
	if len(l) != 2 {
		t.Errorf("wanted list length %d, got %d instead.", 2, len(l))
	}
	if l[1].Task != tasks[2] {
		t.Errorf("wanted %q, got %q instead.", tasks[2], l[1].Task)
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "new task"
	l1.Add(taskName)

	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Errorf("error saving list to file: %s", err)
	}

	if err := l2.Get(tf.Name()); err != nil {
		t.Errorf("error getting list from file: %s", err)
	}

	if l1[0].Task != l2[0].Task {
		t.Errorf("task %v should match %v", l1[0].Task, l2[0].Task)
	}
}

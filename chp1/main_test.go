package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {
	b := bytes.NewBufferString("one two three four\n")

	want := 4

	got := count(b, false)

	if got != want {
		t.Errorf("got %d but want %d", got, want)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("one two three\nfour")

	want := 2

	got := count(b, true)

	if got != want {
		t.Errorf("got %d but want %d", got, want)
	}
}

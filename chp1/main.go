package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// define a bool flag to count lines instead of words
	lines := flag.Bool("l", false, "count lines")
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines))
}

func count(r io.Reader, countLines bool) int {

	// scanner is used to read text from the Reader
	scanner := bufio.NewScanner(r)

	// if countlines flag is not set, we want to split by words
	// default is split by lines so no else statement needed
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0

	// for every word scanned, increment the counter
	for scanner.Scan() {
		wc++
	}

	return wc
}

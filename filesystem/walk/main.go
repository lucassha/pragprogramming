package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext  string
	size int64
	list bool
	del  bool
	wLog io.Writer
}

var (
	f   = os.Stdout
	err error
)

func main() {
	root := flag.String("root", ".", "root directory to start")
	list := flag.Bool("list", false, "list files only")
	ext := flag.String("ext", "", "file extension to filter by")
	size := flag.Int64("size", 0, "minimum file size")
	del := flag.Bool("del", false, "delete files")
	logFile := flag.String("log", "", "log deletes to this file")
	flag.Parse()

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
		wLog: f,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		// if list was explicitly set, don't do anything else
		if cfg.list {
			return listFile(path, out)
		}

		if cfg.del {
			return delFile(path, delLogger)
		}

		// list is the default option is nothing else was set
		return listFile(path, out)
	})
}

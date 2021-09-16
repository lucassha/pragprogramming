package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	// is a directory OR size is too small
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	// ext is not empty AND ext does not equal the val passed in via flag
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	// with filepath you ensure the paths are built in accordance with the
	// operating system where the program is running, making it cross-platform
	targetPath := filepath.Join(destDir, relDir, dest)

	// trying to figure out how this works lol
	fmt.Println("root", root)
	fmt.Println("dir filepath", filepath.Dir(path))
	fmt.Println("rel dir", relDir)
	fmt.Println("target path", targetPath)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(path)

	if _, err = io.Copy(zw, in); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}
	return out.Close()
}

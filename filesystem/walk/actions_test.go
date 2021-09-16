package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, info)

			if f != tc.expected {
				t.Errorf("got '%t' but expected '%t'\n", f, tc.expected)
			}
		})
	}
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		numArchive   int
		numNoArchive int
	}{
		{
			name: "ArchiveExtensionNoMatch",
			cfg: config{
				ext: ".log",
			},
			extNoArchive: ".gz",
			numArchive:   0,
			numNoArchive: 10,
		},
		{
			name: "ArchiveExtensionMatch",
			cfg: config{
				ext: ".log",
			},
			extNoArchive: "",
			numArchive:   10,
			numNoArchive: 0,
		},
		{
			name: "ArchiveExtensionMixed",
			cfg: config{
				ext: ".log",
			},
			extNoArchive: ".gz",
			numArchive:   5,
			numNoArchive: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			// create both the origin dir and the archiving dir
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.numArchive,
				tc.extNoArchive: tc.numNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			expOut := strings.Join(expFiles, "\n")
			got := strings.TrimSpace(buffer.String())

			if got != expOut {
				t.Errorf("got %q but want %q", got, expOut)
			}
		})
	}
}

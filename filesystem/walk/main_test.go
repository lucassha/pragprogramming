package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := ioutil.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// used this piece to see what files are being created in /tmp
	// allfiles, err := ioutil.ReadDir(tempDir)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// for _, f := range allfiles {
	// 	fmt.Println(f.Name())
	// }

	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name: "NoFilter",
			root: "testdata",
			cfg: config{
				ext:  "",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterExtensionMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 0,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 10,
				list: true,
			},
			expected: "testdata/dir.log\n",
		},
		{
			name: "FilterExtensionSizeNoMatch",
			root: "testdata",
			cfg: config{
				ext:  ".log",
				size: 20,
				list: true,
			},
			expected: "",
		},
		{
			name: "FilterExtensionNoMatch",
			root: "testdata",
			cfg: config{
				ext:  ".gz",
				size: 0,
				list: true,
			},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			got := buffer.String()

			if tc.expected != got {
				t.Errorf("got %q but expected %q\n", got, tc.expected)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		numDelete   int
		numNoDelete int
		expected    string
	}{
		{
			name: "DeleteExtensionNoMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			numDelete:   0,
			numNoDelete: 10,
			expected:    "",
		},
		{
			name: "DeleteExtensionMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: "",
			numDelete:   10,
			numNoDelete: 0,
			expected:    "",
		},
		{
			name: "DeleteExtensionMixed",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			numDelete:   5,
			numNoDelete: 5,
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			var logBuffer bytes.Buffer
			tc.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.numDelete,
				tc.extNoDelete: tc.numNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			got := buffer.String()

			if tc.expected != got {
				t.Errorf("got %q but expected %q", got, tc.expected)
			}

			expLogLines := tc.numDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n", expLogLines, len(lines))
			}
		})
	}
}

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	root := filepath.Join("..", "..", "testdata")

	tests := []struct {
		in  string
		golden string
	}{
		{
			"min.jpg",
			"min.jfifcom.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.golden, func(t *testing.T) {
			temp, err := ioutil.TempFile(os.TempDir(), "jfifcom-test-main-"+tt.in)
			if err != nil {
				t.Fatal(err)
			}
			path := temp.Name()
			defer os.Remove(path)
			defer temp.Close()

			src, err := os.Open(filepath.Join(root, tt.in))
			if err != nil {
				t.Fatal(err)
			}
			defer src.Close()

			_, err = io.Copy(temp, src)
			if err != nil {
				t.Fatal(err)
			}
			temp.Close()
			src.Close()

			exit := realMain([]string{path}, strings.NewReader("hello"))
			if exit != 0 {
				t.Fatalf("invalid exit %d", exit)
			}

			compareFiles(t, path, filepath.Join(root, tt.golden))
		})
	}
}

func compareFiles(t *testing.T, path, golden string) {
	want, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	got, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("bytes don't match\ngot:  % x\nwant: % x", got, want) // TODO Better diff
	}
}

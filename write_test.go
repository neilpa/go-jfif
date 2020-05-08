package jfif

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	xio "neilpa.me/go-x/io"
)

func TestAppend(t *testing.T) {
	var tests = []struct {
		name string
		seg Segment
		golden string
	}{
		{
			"min.jpg",
			Segment{Pointer: Pointer{Marker: COM}, Data: []byte("hello")},
			"min.hello.jpg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.golden, func(t *testing.T) {
			path, err := xio.TempFileCopy(filepath.Join("testdata", tt.name), "jfif-test-append-"+tt.golden)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(path)

			err = Append(path, tt.seg)
			if err != nil {
				t.Fatal(err)
			}

			compareFiles(t, path, filepath.Join("testdata", tt.golden))
		})
	}
}

func TestFileQuery(t *testing.T) { // TODO
}

func TestFileUpdate(t *testing.T) {
	var tests = []struct {
		name string
		ref  Pointer
		buf  []byte

		golden string
	}{
		{
			"min.jpg",
			Pointer{Offset: 2, Marker: DQT, Length: 67},
			[]byte{0, // Pq and Tq bytes
				// Arbitrary DQT table for testing
				16, 11, 10, 16, 24, 40, 51, 61,
				12, 12, 14, 19, 26, 58, 60, 55,
				14, 13, 16, 24, 40, 57, 69, 56,
				14, 17, 22, 29, 51, 87, 80, 62,
				18, 22, 37, 56, 68, 109, 103, 77,
				24, 35, 55, 64, 81, 104, 113, 92,
				49, 64, 78, 87, 103, 121, 120, 101,
				72, 92, 95, 98, 112, 100, 103, 99,
			},
			"min.dqt.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.golden, func(t *testing.T) {
			path, err := xio.TempFileCopy(filepath.Join("testdata", tt.name), "jfif-test-update-"+tt.golden)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(path)

			edit, err := Edit(path)
			if err != nil {
				t.Fatal(err)
			}
			defer edit.Close()

			err = edit.Update(tt.ref, tt.buf)
			if err != nil {
				t.Fatal(err)
			}
			edit.Close()

			compareFiles(t, path, filepath.Join("testdata", tt.golden))
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

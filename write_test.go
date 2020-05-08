package jfif

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAppend(t *testing.T) {
	var tests = []struct {
		name string
		ref  Pointer
		buf  []byte

		golden string
	}{
		{
			"min",
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
		t.Run(tt.name, func(t *testing.T) {
			temp, err := ioutil.TempFile(os.TempDir(), "jfif-test-append-"+tt.name)
			if err != nil {
				t.Fatal(err)
			}
			path := temp.Name()
			defer os.Remove(path)
			defer temp.Close()

			fmt.Println("TODO: TestAppend:", path)
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
			"min",
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
		t.Run(tt.name, func(t *testing.T) {
			temp, err := ioutil.TempFile(os.TempDir(), "jfif-test-update-"+tt.name)
			if err != nil {
				t.Fatal(err)
			}
			path := temp.Name()
			defer os.Remove(path)
			defer temp.Close()

			src, err := os.Open(filepath.Join("testdata", tt.name+".jpg"))
			if err != nil {
				t.Fatal(err)
			}
			defer src.Close()

			_, err = io.Copy(temp, src)
			if err != nil {
				t.Fatal(err)
			}
			temp.Close()

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

			want, err := ioutil.ReadFile(filepath.Join("testdata", tt.golden))
			if err != nil {
				t.Fatal(err)
			}
			got, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Error("bytes don't match") // TODO Better diff
			}
		})
	}
}

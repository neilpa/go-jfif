package jfif

import (
	"os"
	"path/filepath"
	"testing"
)

type seg struct {
	marker byte
	size int
}

var tests = []struct{
	path string
	imageSize int
	meta []seg
} {
	{
		path: "lego.jpg",
		imageSize: 216990,
		meta: []seg{
			{0xd8, 0},
			{0xe0, 14},
			{0xe1, 11308},
			{0xe1, 5023},
			{0xdb, 65},
			{0xdb, 65},
			{0xc0, 15},
			{0xc4, 29},
			{0xc4, 79},
			{0xc4, 28},
			{0xc4, 72},
			{0xda, 10},
		},
	},
}

func TestDecodeMetadata(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Fatal(err)
			}
			segments, err := DecodeMetadata(f)
			if err != nil {
				t.Fatal(err)
			}

			if len(segments) != len(tt.meta) {
				t.Fatalf("len: got %d, want %d", len(segments), len(tt.meta))
			}
			for i, s := range segments {
				got := seg{ s.Marker, len(s.Data) }
				want := tt.meta[i]
				if got != want {
					t.Errorf("%d: got %d, want %d", i, got, want)
				}
			}
		})
	}
}

func TestDecodeSegments(t *testing.T) { // TODO
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Fatal(err)
			}
			segments, err := DecodeSegments(f)
			if err != nil {
				t.Fatal(err)
			}

			want := make([]seg, len(tt.meta) + 1)
			copy(want, tt.meta)
			want[len(want)-2].size = tt.imageSize
			want[len(want)-1] = seg { eoiMarker, 0 }

			if len(segments) != len(want) {
				t.Fatalf("len: got %d, want %d", len(segments), len(want))
			}
			for i, s := range segments {
				g := seg{ s.Marker, len(s.Data) }
				w := want[i]
				if g != w {
					t.Errorf("%d: got %d, want %d", i, g, w)
				}
			}
		})
	}
}

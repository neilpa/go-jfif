package jfif

import (
	"os"
	"path/filepath"
	"testing"
)

type seg struct {
	marker Marker
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
			{SOI, 0},
			{APP0, 14},
			{APP1, 11308},
			{APP1, 5023},
			{DQT, 65},
			{DQT, 65},
			{SOF0, 15},
			{DHT, 29},
			{DHT, 79},
			{DHT, 28},
			{DHT, 72},
			{SOS, 10 },
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
			want[len(want)-1] = seg { EOI, 0 }

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

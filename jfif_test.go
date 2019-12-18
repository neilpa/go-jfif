package jfif

import (
	"os"
	"path/filepath"
	"testing"
)

type seg struct {
	marker Marker
	size   int
	offset int64
}

var tests = []struct {
	path      string
	imageSize int
	meta      []seg
}{
	{
		path:      "lego.jpg",
		imageSize: 216990,
		meta: []seg{
			{marker: SOI, size: 0},
			{marker: APP0, size: 14},
			{marker: APP1, size: 11308},
			{marker: APP1, size: 5023},
			{marker: DQT, size: 65},
			{marker: DQT, size: 65},
			{marker: SOF0, size: 15},
			{marker: DHT, size: 29},
			{marker: DHT, size: 79},
			{marker: DHT, size: 28},
			{marker: DHT, size: 72},
			{marker: SOS, size: 10},
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

			verifySegments(t, segments, tt.meta)
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

			want := make([]seg, len(tt.meta)+1)
			copy(want, tt.meta)
			want[len(want)-2].size = tt.imageSize // TODO + ImageSize?
			want[len(want)-1] = seg{EOI, 0, -1}

			verifySegments(t, segments, want)
		})
	}
}


func verifySegments(t *testing.T, segments []Segment, want []seg) {
	if len(segments) != len(want) {
		t.Fatalf("len: got %d, want %d", len(segments), len(want))
	}

	var offset int64
	for i, s := range segments {
		g := seg{s.Marker, len(s.Data), s.Offset}
		w := want[i]
		w.offset = offset
		if g != w {
			t.Fatalf("%d: got %d, want %d", i, g, w)
		}
		// 0xff and marker
		offset += 2
		if s.Data != nil {
			// 2-byte length and data
			offset += 2 + int64(len(s.Data))
		}
	}
}

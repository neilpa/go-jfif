package jfif

import (
	"os"
	"path/filepath"
	"testing"

	"fmt"
)

var tests = []struct {
	path string
	meta []SegmentRef
}{
	{
		path:      "lego.jpg",
		meta: []SegmentRef{
			{Marker: SOI, Size: 2},
			{Marker: APP0, Size: 18},
			{Marker: APP1, Size: 11312},
			{Marker: APP1, Size: 5027},
			{Marker: DQT, Size: 69},
			{Marker: DQT, Size: 69},
			{Marker: SOF0, Size: 19},
			{Marker: DHT, Size: 33},
			{Marker: DHT, Size: 83},
			{Marker: DHT, Size: 32},
			{Marker: DHT, Size: 76},
			{Marker: SOS, Size: 14},
		},
	},
}

func TestDecodeSegments(t *testing.T) {
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
			refs := make([]SegmentRef, len(segments))
			for i, seg := range segments {
				refs[i] = seg.SegmentRef
			}
			verifySegments(t, refs, tt.meta)
		})
	}
}

func TestDecodeMetadata(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Fatal(err)
			}
			refs, err := DecodeMetadata(f)
			if err != nil {
				t.Fatal(err)
			}
			verifySegments(t, refs, tt.meta)
		})
	}
}

func TestEncodeSegments(t *testing.T) { // TODO
}

func verifySegments(t *testing.T, got, want []SegmentRef) {
	if len(got) != len(want) {
		fmt.Println(got)
		fmt.Println(want)
		t.Fatalf("len: got %d, want %d", len(got), len(want))
	}

	var offset int64
	for i, g := range got {
		w := want[i]
		w.Offset = offset
		if g != w {
			t.Fatalf("%d: got %#v, want %#v", i, g, w)
		}
		offset += int64(g.Size)
	}
}

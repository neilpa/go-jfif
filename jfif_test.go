package jfif

import (
	"os"
	"path/filepath"
	"testing"
)

var tests = []struct {
	path string
	refs []SegmentP
}{
	{
		path: "min.jpg",
		refs: []SegmentP{
			{0, SOI, 0},
			{2, DQT, 67},
			{71, Marker(0xC9), 11},
			{84, Marker(0xCC), 6},
			{92, SOS, 8},
		},
	},
	{
		path: "lego.jpg",
		refs: []SegmentP{
			{0, SOI, 0},
			{2, APP0, 16},
			{20, APP1, 11310},
			{11332, APP1, 5025},
			{16359, DQT, 67},
			{16428, DQT, 67},
			{16497, SOF0, 17},
			{16516, DHT, 31},
			{16549, DHT, 81},
			{16632, DHT, 30},
			{16664, DHT, 74},
			{16740, SOS, 12},
		},
	},
}

func TestScanSegments(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Fatal(err)
			}
			refs, err := ScanSegments(f)
			if err != nil {
				t.Fatal(err)
			}

			verifySegments(t, refs, tt.refs)
		})
	}
}

func TestDecodeSegments(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", tt.path))
			if err != nil {
				t.Fatal(err)
			}
			segs, err := DecodeSegments(f)
			if err != nil {
				t.Fatal(err)
			}

			refs := make([]SegmentP, len(segs))
			for i, s := range segs {
				refs[i] = s.SegmentP
				if s.Length > 0 && len(s.Data)+2 != int(s.Length) {
					t.Errorf("data %d: got %d, want %d", i, len(s.Data), s.Length-2)
				}
			}
			verifySegments(t, refs, tt.refs)
		})
	}
}

func TestEncodeSegment(t *testing.T) { // TODO
}

func verifySegments(t *testing.T, got, want []SegmentP) {
	if len(got) != len(want) {
		t.Errorf("len: got %d, want %d", len(got), len(want))
		return
	}

	for i, w := range want {
		g := got[i]
		if g != w {
			t.Errorf("seg %d: got %d, want %d", i, g, w)
		}
	}
}

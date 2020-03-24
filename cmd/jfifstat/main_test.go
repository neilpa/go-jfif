// jfifstat prints segment markers, sizes, and optional APPN signatures from
// JPEG files. Stops after the Start of Stream (SOS) segment.
package main

import (
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	tests := []struct {
		in string
		out []string
	} {
		{
			"../../testdata/lego.jpg",
			[]string{
				"SOI	0",
				"APP0	14	JFIF",
				"APP1	11308	Exif",
				"APP1	5023	http://ns.adobe.com/xap/1.0/",
				"DQT	65",
				"DQT	65",
				"SOF0	15",
				"DHT	29",
				"DHT	79",
				"DHT	28",
				"DHT	72",
				"SOS	10",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			var stdout strings.Builder
			exit := realMain(strings.Split(tt.in, " "), &stdout)
			if exit != 0 {
				t.Fatalf("exit: %d", exit)
			}
			//got := strings.TrimSpace(stdout.String())
			//want := strings.Join(tt.out, "\n")
			//if got != want {
			//	t.Errorf("got:\n%s\nwant:\n%s\n", got, want)
			//}

			got := strings.Split(strings.TrimSpace(stdout.String()), "\n")
			if len(got) != len(tt.out) {
				t.Fatalf("len: got %d want %d", len(got), len(tt.out))
			}
			for i, line := range got {
				if line != tt.out[i] {
					t.Errorf("line %d: got %q want %q", i, line, tt.out[i])
				}
			}
		})
	}
}

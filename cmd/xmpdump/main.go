package main

import (
	"bytes"
	"log"
	"os"
	"strings"

	"neilpa.me/go-jfif"

	"trimmer.io/go-xmp/xmp"
)

const (
	sigXMP = "http://ns.adobe.com/xap/1.0/\x00"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	segs, err := jfif.DecodeMetadata(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range segs {
		if s.Marker == jfif.APP1 {
			// TODO Support ExtendedXMP with merging
			// https://stackoverflow.com/questions/23253281/reading-jpg-files-xmp-metadata
			if strings.HasPrefix(string(s.Data), sigXMP) {
				packets, err := xmp.ScanPackets(bytes.NewReader(s.Data))
				if err != nil {
					log.Fatal(err)
				}

				for _, p := range packets {
					doc, err := xmp.Read(bytes.NewReader(p))
					if err != nil {
						log.Fatal(err)
					}
					doc.Dump()
				}
			}
		}
	}
}

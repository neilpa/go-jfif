package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"neilpa.me/go-jfif"
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
		sig, payload, _ := s.AppPayload()
		switch sig {
		case jfif.SigXMP:
			fmt.Println(string(payload))
		case jfif.SigExtendedXMP:
			var h header
			r := bytes.NewReader(payload)
			err = binary.Read(r, binary.BigEndian, &h)
			if err != nil {
				log.Fatal(err)
			}
			// TODO Buffer if different extendedd XMP segments are
			// interleaved or they are in the wrong order
			fmt.Print(string(payload[binary.Size(h):]))
		}
	}
}

// Each chunk is written into the JPEG file within a separate APP1 marker
// segment. Each ExtendedXMP marker segment contains:
//
//  A null-terminated signature string "http://ns.adobe.com/xmp/extension/".
//  A 128-bit GUID stored as a 32-byte ASCII hex string, capital A-F, no
//    null termination. The GUID is a 128-bit MD5 digest of the full
//    ExtendedXMP serialization.
//  The full length of the ExtendedXMP serialization as a 32-bit unsigned
//    integer.
//  The offset of this portion as a 32-bit unsigned integer.
//  The portion of the ExtendedXMP.

type header struct {
	MD5Hash    [32]byte
	FullLength uint32
	Offset     uint32
}

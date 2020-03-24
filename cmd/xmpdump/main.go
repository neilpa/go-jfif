// xmpdump extracts and prints XMP data from APP1 segments from JPEG/JFIF files.
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"neilpa.me/go-jfif"
)

func main() {
	os.Exit(realMain(os.Args[1:]))
}

func realMain(args []string) int {
	flag.Usage = printUsage
	flag.CommandLine.Parse(args)

	if flag.NArg() == 0 {
		return usageError("no files specified")
	}

	for _, arg := range flag.Args() {
		f, err := os.Open(arg)
		if err != nil {
			return fatal(err.Error())
		}

		segs, err := jfif.DecodeMetadata(f)
		if err != nil {
			return fatal(err.Error())
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
					return fatal(err.Error())
				}
				// TODO Buffer if different extendedd XMP segments are
				// interleaved or they are in the wrong order
				fmt.Print(string(payload[binary.Size(h):]))
			}
		}
	}
	return 0
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

func fatal(format string, args ...interface{}) int {
	format = os.Args[0] + ": " + format + "\n"
	fmt.Fprintf(os.Stderr, format, args...)
	return 1
}

func usageError(msg string) int {
	fmt.Fprintln(os.Stderr, msg)
	printUsage()
	return 2
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s file [file...]

  xmpdump extracts and prints XMP data from APP1 segments from JPEG/JFIF files.
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

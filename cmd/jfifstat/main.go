// jfifstat prints segment markers, sizes, and optional APPN signatures from
// JPEG files. Stops after the Start of Stream (SOS) segment.
package main

import (
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
			sig, _, _ := s.AppPayload()
			fmt.Printf("%s\t%d\t%s\n", s.Marker, len(s.Data), sig)
		}
	}
	return 0
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
	fmt.Fprintf(os.Stderr, `Usage: %s jpeg [jpeg...]

  jfifstat prints segment markers, sizes, and optional APPN signatures from
  JPEG files. Stops after the Start of Stream (SOS) segment.
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

// jfifstat prints segment markers, sizes, and optional APPN signatures from
// JPEG files. Stops after the Start of Stream (SOS) segment. Segment lines
// are prefixed with the filename when multiple files are specified.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"neilpa.me/go-jfif"
)

func main() {
	os.Exit(realMain(os.Args[1:], os.Stdout))
}

func realMain(args []string, stdout io.Writer) int {
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

		segs, err := jfif.DecodeSegments(f)
		if err != nil {
			return fatal(err.Error())
		}

		for _, s := range segs {
			sig, _, _ := s.AppPayload()
			sig = jfif.CleanSig(sig)
			if flag.NArg() > 1 {
				if sig != "" {
					fmt.Fprintf(stdout, "%s\t%s\t%d\t%s\n", arg, s.Marker, len(s.Data), sig)
				} else {
					fmt.Fprintf(stdout, "%s\t%s\t%d\n", arg, s.Marker, len(s.Data))
				}
			} else if sig != "" {
				fmt.Fprintf(stdout, "%s\t%d\t%s\n", s.Marker, len(s.Data), sig)
			} else {
				fmt.Fprintf(stdout, "%s\t%d\n", s.Marker, len(s.Data))
			}
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
  JPEG files. Stops after the Start of Stream (SOS) segment. Segment lines
  are prefixed with the filename when multiple files are specified.
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

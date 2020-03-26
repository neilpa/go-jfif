// jfifsig guesses signatures of APPn segments in a JPEG files. If it
// matches a known signature it prints that directly. Otherwise, it
// tries to figure it out by looking for non-printable characters
// up to some limit.
// TODO Should just make this an option for jfifstat
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

		segs, err := jfif.DecodeSegments(f) // TODO Scan Segments
		if err != nil {
			return fatal(err.Error())
		}

		for _, s := range segs {
			if s.Marker < jfif.APP0 || s.Marker > jfif.APP15 {
				continue
			}

			sig, _, _ := s.AppPayload()
			sig = jfif.CleanSig(sig)
			if len(sig) > 0 {
				if flag.NArg() > 1 {
					fmt.Printf("%s\tmatch\t%s\t%s\n", arg, s.Marker, sig)
				} else {
					fmt.Printf("match\t%s\t%s\n", s.Marker, sig)
				}
				continue
			}

			// Guessing game
			// TODO Actually look null (\x00) bytes as terminators
			limit := 20
			if len(s.Data) < limit {
				limit = len(s.Data)
			}
			guess := string(s.Data[:limit])
			if flag.NArg() > 1 {
				fmt.Fprintf(stdout, "%s\tguess\t%s\t%q\n", arg, s.Marker, guess)
			} else {
				fmt.Fprintf(stdout, "guess\t%s\t%s\n", s.Marker, guess)
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

  jfifsig guesses signatures of APPn segments in a JPEG files. If it
  matches a known signature it prints that directly. Otherwise, it
  tries to figure it out by looking for non-printable characters
  up to some limit.
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

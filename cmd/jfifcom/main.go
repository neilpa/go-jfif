// jfifcom embeds a new comment segment with the data from stdin
// before the start of stream (SOS) segment.
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"neilpa.me/go-jfif"
)

func main() {
	os.Exit(realMain(os.Args[1:], os.Stdin))
}

func realMain(args []string, stdin io.Reader) int {
	flag.Usage = printUsage
	flag.CommandLine.Parse(args)

	if flag.NArg() == 0 {
		return usageError("no files specified")
	}

	buf, err := ioutil.ReadAll(stdin)
	if err != nil {
		return fatal(err.Error())
	}

	for _, arg := range flag.Args() {
		//fmt.Println("embeddding", arg, "buf", buf)

		err = jfif.Add(arg, jfif.COM, buf)
		if err != nil {
			return fatal(err.Error()) // todo: continue writing the other files?
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
	fmt.Fprintf(os.Stderr, `Usage: %s jpeg [jpeg...] < ...

  jfifcomment embeds a new comment segment with the data from stdin
  before the start of stream (SOS) segment.
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
}

package main

import (
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

	segs, err := jfif.DecodeSegments(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range segs {
		fmt.Printf("%s\t%d\n", s.Marker, len(s.Data))
	}
}

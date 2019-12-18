package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"neilpa.me/go-jfif"
)

const (
	sigXMP = "http://ns.adobe.com/xap/1.0/\x00"
	sigExtendedXMP = "http://ns.adobe.com/xmp/extension/\x00"
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
			data := string(s.Data)
			if strings.HasPrefix(data, sigXMP) {
				fmt.Println(data)
			} else if strings.HasPrefix(data, sigExtendedXMP) {
				// TODO Merge data that spans multiple segments
				fmt.Println(data)
			}
		}
	}
}

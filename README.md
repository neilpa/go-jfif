# go-jfif

[![CI](https://github.com/neilpa/go-jfif/workflows/CI/badge.svg)](https://github.com/neilpa/go-jfif/actions/)

Basic reading of [JPEG File Interchange Format (JFIF)][wiki-jfif] segments

[wiki-jfif]: https://en.wikipedia.org/wiki/JPEG_File_Interchange_Format#File_format_structure

## Install

Use `go get` to install the module

```
go get neilpa.me/go-jfif
```

Or use it and rely on go modules to "do the right thing".

## Usage

Reading JFIF segments from an existing JPEG.

```go
package main

import (
	"fmt"
	"log"
	"os"

	"neilpa.me/go-jfif"
)

func main() {
	f, err := os.Open("path/to/file.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// See also jfif.ScanSegments which doesn't read the segment payload.
	// This is used to detect the segment "signatures" for some APPn segments.
	segs, err := jfif.DecodeSegments(f)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range segs {
		sig, _, _ := s.AppPayload()
		sig = jfif.CleanSig(sig)
		fmt.Printf("%d\t%s\t%s\n", s.Length, s.Marker, sig)
	}
}
```

Appending new segments to a JPEG. (Under the hood this edits a copy of the file before renaming the copy
to finalize the updates).

```go
package main

import (
	"log"

	"neilpa.me/go-jfif"
)

func main() {
	err := jfif.Add("path/to/file.jpg", jfif.COM, []byte("adding a comment to this file"))
	if err != nil {
		log.Fatal(err)
	}
}
```

## Tools

There are a few simple CLI tools include for manipulating JPEG/JFIF files that exercise functionality
exposed by this module.

### jfifcom

Append free-form comment (`COM`) segments to an existing JPEG file.

### jfifstat

Prints segment markers, sizes, and optional `APPN` signatures from a JPEG until the image stream.

### xmpdump

Extracts and prints [XMP](https://www.adobe.com/products/xmp.html) data from `APP1` segments from JPEG file(s).

## License

[MIT](/LICENSE)

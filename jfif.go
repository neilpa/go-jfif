// Package jfif supports the basics of reading segments from JPEG files.
//
// https://en.wikipedia.org/wiki/JPEG#Syntax_and_structure
package jfif // import "neilpa.me/go-jfif"

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	soiMarker = 0xd8 // Start Of Image.
	eoiMarker = 0xd9 // End Of Image.
	sosMarker = 0xda // Start Of Stream/Scan.
)

var (
	// ErrInvalid means the reader did not begin with a start of image
	// marker.
	ErrInvalid = errors.New("Invalid JPEG")

	// ErrShortSegment means a segment length was < 2 bytes.
	ErrShortSegment = errors.New("Short segment")
)

// Segment identifies a part of a JPEG file and associatedd raw data.
type Segment struct {
	Marker byte
	Data []byte
}

// DecodeMetadata reads segments until the start of stream (SOS) marker is read,
// or an error is encountered, including EOF. This will read the SOS segment but
// not the subsequent entropy-coded image data.
// TODO Should this return "io.ErrUnexpectedEOF" when EOF is seen before SOS?
func DecodeMetadata(r io.Reader) ([]Segment, error) {
	var magic [2]byte
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic[0] != 0xff || magic[1] != soiMarker {
		return nil, ErrInvalid
	}

	// This behavior matches that of image/jpeg.decode
	// https://golang.org/src/image/jpeg/reader.go?s=22312:22357#L526
	segments := []Segment{{magic[1], nil}}
	for {
		var buf [2]byte
		err = binary.Read(r, binary.BigEndian, &buf)
		if err != nil {
			return segments, err
		}
		sentinel, marker := buf[0], buf[1]

		for sentinel != 0xff {
			// Technically a format error but mimics go's stdlib which is
			// itself matching the behavor of libjpeg.
			sentinel = marker
			marker, err = readByte(r)
			if err != nil {
				return segments, err
			}
		}

		if marker == 0 {
			// Byte Stuffing, e.g. "Extraneous Data"
			// TODO Does this actually matter if reading to EOI once the
			// SOS marker is seen? If so, should these be included?
			continue
		}

		for marker == 0xff {
			// Eat fill bytes that may precede a marker
			// TODO Does this actually matter if reading to EOI once the
			// SOS marker is seen?
			marker, err = readByte(r)
			if err != nil {
				return segments, err
			}
		}

		var length uint16 // TODO Is this an int16?
		if err = binary.Read(r, binary.BigEndian, &length); err != nil {
			return segments, err
		}
		if length < 2 {
			return segments, ErrShortSegment
		}

		// Length includes the 2 bytes for itself
		s := Segment{marker, make([]byte, int(length) - 2) }
		if err = binary.Read(r, binary.BigEndian, &s.Data); err != nil {
			return segments, err
		}
		segments = append(segments, s)

		if marker == sosMarker {
			break
		}

	}

	return segments, nil
}

// DecodeSegments reads segments until the end of image (EOI) marker is read, or an
// error is encountered, including EOF.
// TODO Should this return "io.ErrUnexpectedEOF" when io.EOF is seen before EOI?
func DecodeSegments(r io.Reader) ([]Segment, error) {
	segments, err := DecodeMetadata(r)
	if err != nil {
		return segments, err
	}
	return segments, errors.New("TODO: Read until EOI")
}

func readByte(r io.Reader) (b byte, err error) {
	err = binary.Read(r, binary.BigEndian, &b)
	return
}

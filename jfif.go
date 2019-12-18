// Package jfif supports the basics of reading segments from JPEG files.
//
// https://en.wikipedia.org/wiki/JPEG#Syntax_and_structure
package jfif // import "neilpa.me/go-jfif"

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

var (
	// ErrInvalid means the reader did not begin with a start of image
	// marker.
	ErrInvalid = errors.New("Invalid JPEG")

	// ErrShortSegment means a segment length was < 2 bytes.
	ErrShortSegment = errors.New("Short segment")
)

// Segment represents a distinct region of a JPEG file.
type Segment struct {
	// Marker identifies the type of segment.
	Marker
	// Data is the raw bytes of a segment, excluding the initial 4 bytes
	// (e.g. 0xff, marker, and 2-byte length). For segments lacking a
	// length, this will be nil.
	Data []byte
	// Offset is the address of the 0xff byte that started this segment that
	// is then followed by the marker.
	Offset int64
}

// DecodeMetadata reads segments until the start of stream (SOS) marker is
// read, or an error is encountered, including EOF. This will read the SOS
// segment and its payload but not the subsequent entropy-coded image data.
func DecodeMetadata(r io.Reader) ([]Segment, error) {
	counter, ok := r.(*countReader)
	if !ok {
		counter = &countReader{reader: r}
	}
	r = counter

	var magic [2]byte
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic[0] != 0xff || magic[1] != byte(SOI) {
		return nil, ErrInvalid
	}

	// This behavior matches that of image/jpeg.decode
	// https://golang.org/src/image/jpeg/reader.go?s=22312:22357#L526
	segments := []Segment{{Marker: Marker(magic[1])}}
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

		// Set the offset to the 0xff byte preceding the marker
		s := Segment{Marker: Marker(marker), Offset: counter.count - 2}

		var length uint16 // TODO Is this an int16?
		if err = binary.Read(r, binary.BigEndian, &length); err != nil {
			return segments, err
		}
		if length < 2 {
			return segments, ErrShortSegment
		}

		// Length includes the 2 bytes for itself
		s.Data = make([]byte, int(length)-2)
		if err = binary.Read(r, binary.BigEndian, &s.Data); err != nil {
			return segments, err
		}
		segments = append(segments, s)

		if marker == SOS {
			break
		}
	}

	return segments, nil
}

// DecodeSegments reads segments through an end of image (EOI) marker, or
// an error is encountered, including EOF. The entropy-coded image data
// _should_ be the penulatimate segment, following the SOS segment and
// with the fabricated XXX marker as an identifer.
func DecodeSegments(r io.Reader) ([]Segment, error) {
	counter, ok := r.(*countReader)
	if !ok {
		counter = &countReader{reader: r}
	}
	r = counter

	segments, err := DecodeMetadata(r)
	if err != nil {
		return segments, err
	}

	seg := Segment{Marker: XXX, Offset: counter.count}
	b := bufio.NewReader(r)
	for {
		data, err := b.ReadBytes(0xff)
		if err != nil {
			return segments, err
		}
		data = data[:len(data)-1]
		if seg.Data == nil {
			seg.Data = data
		} else {
			seg.Data = append(seg.Data, data...)
		}

		marker, err := b.ReadByte()
		if err != nil {
			return segments, err
		}
		if marker == EOI {
			eoi := Segment{Marker(marker), nil, counter.count - 2}
			segments = append(segments, seg, eoi)
			break
		}
		// Add back the sentinal and marker and continue
		seg.Data = append(seg.Data, 0xff, marker)
	}

	return segments, nil
}

func readByte(r io.Reader) (b byte, err error) {
	err = binary.Read(r, binary.BigEndian, &b)
	return
}

type countReader struct {
	reader io.Reader
	count  int64
}

func (c *countReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	c.count += int64(n)
	return
}

// Package jfif supports simple reading and writing of segments from a JPEG
// file.
//
// https://en.wikipedia.org/wiki/JPEG#Syntax_and_structure
package jfif // import "neilpa.me/go-jfif"

import (
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

	// ErrWrongMarker means a segment method was called where the marker
	// didn't match an expected type.
	ErrWrongMarker = errors.New("Wrong marker")

	// ErrUnknownApp means an APPn segment has an unrecognized signature.
	ErrUnknownApp = errors.New("Unknown APPn segment")
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

// AppPayload extracts a recognized signature and payload bytes. Otherwise
// returns an erorr for non APPn segments or if the signature is unknown.
func (s Segment) AppPayload() (string, []byte, error) {
	if s.Marker < APP0 || s.Marker > APP15 {
		return "", nil, ErrWrongMarker
	}
	for _, sig := range appnSigs[int(s.Marker-APP0)] {
		if sig == string(s.Data[:len(sig)]) {
			return sig, s.Data[len(sig):], nil
		}
	}
	return "", nil, ErrUnknownApp
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

// EncodeSegment writes the given segment.
func EncodeSegment(w io.Writer, seg Segment) error {
	// Everything needs the 0xff, marker and potential payload
	_, err := w.Write([]byte{0xff, byte(seg.Marker)})
	if err != nil || seg.Data == nil {
		return err
	}
	// Payload size includes it's own 2-bytes
	err = binary.Write(w, binary.BigEndian, uint16(len(seg.Data))+2)
	if err != nil {
		return err
	}
	_, err = w.Write(seg.Data)
	return err
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

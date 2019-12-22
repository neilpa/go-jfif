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

// DecodeSegments reads segments until the start of stream (SOS) marker is
// read, or an error is encountered, including EOF. This will read the SOS
// segment and its payload but not the subsequent entropy-coded image data.
func DecodeSegments(r io.Reader) ([]Segment, error) {
	dec, err := newDecoder(r)
	if err != nil {
		return nil, err
	}
	segments := []Segment{emptySegment(SOI, 0)}
	for {
		marker, err := nextMarker(dec)
		if err != nil {
			return segments, err
		}
		s, err := decodeSegment(dec, marker)
		if err != nil {
			return segments, err
		}
		segments = append(segments, s)
		if marker == SOS {
			break
		}
	}
	return segments, nil
}

// DecodeMetadata scans segments until the start of stream (SOS) marker,
// or an error is encountered, including EOF. This will seek past the SOS
// segment to the start of the image data.
func DecodeMetadata(r io.ReadSeeker) ([]SegmentRef, error) {
	dec, err := newDecoder(r)
	if err != nil {
		return nil, err
	}
	refs := []SegmentRef{emptyRef(SOI, 0)}
	for {
		marker, err := nextMarker(dec)
		if err != nil {
			return refs, err
		}
		ref, err := decodeRef(dec, marker)
		if err != nil {
			return refs, err
		}
		refs = append(refs, ref)
		err = dec.Skip(int64(ref.Size - 4))
		if err != nil {
			return refs, err
		}
		if marker == SOS {
			break
		}
	}
	return refs, nil
}

// EncodeSegments writes the given segments.
func EncodeSegments(w io.Writer, segments []Segment) error {
	for _, s := range segments {
		err := s.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}

// newDecoder validates the header and returns a byte counting reader.
func newDecoder(r io.Reader) (*countReader, error) {
	var magic [2]byte
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return nil, err
	}
	if magic[0] != 0xff || magic[1] != byte(SOI) {
		return nil, ErrInvalid
	}
	return &countReader{r, r.(io.Seeker), 2}, nil
}

// nextMarker finds the marker at the beginning of the next segment.
func nextMarker(dec *countReader) (Marker, error) {
	// This behavior matches that of image/jpeg.decode
	// https://golang.org/src/image/jpeg/reader.go?s=22312:22357#L526
	var buf [2]byte
	for {
		err := binary.Read(dec, binary.BigEndian, &buf)
		if err != nil {
			return 0, err
		}
		sentinel, marker := buf[0], buf[1]

		for sentinel != 0xff {
			// Technically a format error but mimics go's stdlib which is
			// itself matching the behavior of libjpeg.
			sentinel = marker
			marker, err = readByte(dec)
			if err != nil {
				return 0, err
			}
		}

		if marker == 0 {
			// Byte Stuffing, e.g. "Extraneous Data"
			continue
		}

		for marker == 0xff {
			// Eat fill bytes that may precede a marker
			marker, err = readByte(dec)
			if err != nil {
				return 0, err
			}
		}

		return Marker(marker), nil
	}
}

func readByte(r io.Reader) (b byte, err error) {
	err = binary.Read(r, binary.BigEndian, &b)
	return
}

type countReader struct {
	reader io.Reader
	seeker io.Seeker
	count  int64
}

func (c *countReader) Read(p []byte) (n int, err error) {
	n, err = c.reader.Read(p)
	c.count += int64(n)
	return
}

func (c *countReader) Skip(n int64) (err error) {
	c.count, err = c.seeker.Seek(n, io.SeekCurrent)
	return
}

// Package jfif supports simple reading and writing of segments from a JPEG
// file.
//
// https://en.wikipedia.org/wiki/JPEG#Syntax_and_structure
package jfif // import "neilpa.me/go-jfif"

import (
	"encoding/binary"
	"errors"
	"io"
	"strings"

	xio "neilpa.me/go-x/io"
)

var (
	// ErrInvalidJPEG means the reader did not begin with a start of image
	// marker.
	ErrInvalidJPEG = errors.New("Invalid JPEG")

	// ErrShortSegment means a segment length was < 2 bytes.
	ErrShortSegment = errors.New("Short segment")

	// ErrWrongMarker means a segment method was called where the marker
	// didn't match an expected type.
	ErrWrongMarker = errors.New("Wrong marker")

	// ErrUnknownApp means an APPn segment has an unrecognized signature.
	ErrUnknownApp = errors.New("Unknown APPn segment")

	// ErrUnseekableReader means a Seek was attempted from the start or end
	// of an io.Reader that only supports streaming.
	ErrUnseekableReader = errors.New("Unseekable reader")
)

// SegmentP represents a "pointer" to a distinct region of a JPEG file.
type SegmentP struct {
	// Offset is the address of the 0xff byte that started this segment that
	// is then followed by the marker.
	Offset int64
	// Marker identifies the type of segment.
	Marker
	// Length is the 2-byte segment size after the Marker. Note it's
	// inclusive of the bytes to store it, e.g. len(Data) = Length-2.
	Length uint16
}

// Segment represents a distinct region of a JPEG file.
type Segment struct {
	// SegmentP embeds the positional information of the segment.
	SegmentP
	// Data is the raw bytes of a segment, excluding the initial 4 bytes
	// (e.g. 0xff, marker, and 2-byte length). For segments lacking a
	// length it will be nil.
	Data []byte
}

// AppPayload extracts a recognized signature and payload bytes. Otherwise
// returns an erorr for non APPn segments or if the signature is unknown.
// Note many known signatures include non-printable suffixes like '\0', use
// CleanSig to strip these.
func (s Segment) AppPayload() (string, []byte, error) {
	if s.Marker < APP0 || s.Marker > APP15 {
		return "", nil, ErrWrongMarker
	}
	for _, sig := range appnSigs[int(s.Marker-APP0)] {
		if strings.HasPrefix(string(s.Data), sig) {
			return sig, s.Data[len(sig):], nil
		}
	}
	return "", nil, ErrUnknownApp
}

// ScanSegments finds segment markers until the start of stream (SOS)
// marker is read, or an error is encountered, including EOF.
func ScanSegments(r io.Reader) ([]SegmentP, error) {
	var segs []SegmentP
	err := readSegments(r, func(r io.ReadSeeker, sp SegmentP) error {
		if sp.Length > 0 {
			// Simply skip past the length of the segment
			if _, err := r.Seek(int64(sp.Length)-2, io.SeekCurrent); err != nil {
				return err
			}
		}
		segs = append(segs, sp)
		return nil
	})
	return segs, err
}

// DecodeSegments reads segments and payloads through the start of stream
// (SOS) marker, or until an error is encountered, including EOF. On success
// the reader will be positioned at the start of the entropy-coded image
// data.
func DecodeSegments(r io.Reader) ([]Segment, error) {
	var segs []Segment
	err := readSegments(r, func(r io.ReadSeeker, sp SegmentP) error {
		s := Segment{SegmentP: sp}
		if s.Length > 0 {
			// Length includes the 2 bytes for itself
			s.Data = make([]byte, int(s.Length)-2)
			if _, err := io.ReadFull(r, s.Data); err != nil {
				return err
			}
		}
		segs = append(segs, s)
		return nil
	})
	return segs, err
}

// readSegments scans for segment start markers and calculates the length.
// The provided function is then called with each segment for processing
// the payload data. This function must advance the reader to the end of the
// given segment for the next read.
// TODO Could forego that requirement given the use of xio.TrackingReader
func readSegments(r io.Reader, fn func(io.ReadSeeker, SegmentP) error) error {
	tr, ok := r.(*xio.TrackingReader)
	if !ok {
		tr = xio.NewTrackingReader(r)
	}
	r = tr

	var magic [2]byte
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return err
	}
	if magic[0] != 0xff || magic[1] != byte(SOI) {
		return ErrInvalidJPEG
	}

	err = fn(tr, SegmentP{Marker: Marker(magic[1])})
	if err != nil {
		return err
	}

	// This behavior matches that of image/jpeg.decode
	// https://golang.org/src/image/jpeg/reader.go?s=22312:22357#L526
	for {
		var buf [2]byte
		err = binary.Read(r, binary.BigEndian, &buf)
		if err != nil {
			return err
		}
		sentinel, marker := buf[0], buf[1]

		for sentinel != 0xff {
			// Technically a format error but mimics go's stdlib which is
			// itself matching the behavior of libjpeg.
			sentinel = marker
			marker, err = readByte(r)
			if err != nil {
				return err
			}
		}
		if marker == 0 {
			// Byte Stuffing, e.g. "Extraneous Data"
			continue
		}
		for marker == 0xff {
			// Eat fill bytes that may precede a marker
			marker, err = readByte(r)
			if err != nil {
				return err
			}
		}

		// Set the offset to the 0xff byte preceding the marker
		s := SegmentP{Marker: Marker(marker), Offset: tr.Offset() - 2}

		// TODO Are there expected zero-length markers that can be skipped
		if err = binary.Read(r, binary.BigEndian, &s.Length); err != nil {
			return err
		}
		if s.Length < 2 {
			return ErrShortSegment
		}
		if err = fn(tr, s); err != nil {
			return err
		}
		if marker == byte(SOS) {
			break
		}
	}

	return nil
}

// EncodeSegment writes the given segment.
func EncodeSegment(w io.Writer, seg Segment) error {
	// Everything else needs the 0xff, marker and potential payload
	_, err := w.Write([]byte{0xff, byte(seg.Marker)})
	if err != nil || seg.Data == nil {
		return err
	}
	// Payload size includes it's own 2-bytes
	// TODO Validate the lenght of Data here?
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

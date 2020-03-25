// Package jfif supports simple reading and writing of segments from a JPEG
// file.
//
// https://en.wikipedia.org/wiki/JPEG#Syntax_and_structure
package jfif // import "neilpa.me/go-jfif"

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"strings"
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
	err := readSegments(r, func(r *positionalReader, sp SegmentP) error {
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
	err := readSegments(r, func(r *positionalReader, sp SegmentP) error {
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
// TODO Could forego that requirement given the use of positionalReader
func readSegments(r io.Reader, fn func(*positionalReader, SegmentP) error) error {
	pr, ok := r.(*positionalReader)
	if !ok {
		pr = &positionalReader{reader: r}
	}
	r = pr

	var magic [2]byte
	err := binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return err
	}
	if magic[0] != 0xff || magic[1] != byte(SOI) {
		return ErrInvalid
	}

	err = fn(pr, SegmentP{Marker: Marker(magic[1])})
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
			// itself matching the behavor of libjpeg.
			sentinel = marker
			marker, err = readByte(r)
			if err != nil {
				return err
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
				return err
			}
		}

		// Set the offset to the 0xff byte preceding the marker
		s := SegmentP{Marker: Marker(marker), Offset: pr.pos - 2}

		// TODO Are there expected zero-length markers that can be skipped
		if err = binary.Read(r, binary.BigEndian, &s.Length); err != nil {
			return err
		}
		if s.Length < 2 {
			return ErrShortSegment
		}
		if err = fn(pr, s); err != nil {
			return err
		}
		if marker == SOS {
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

// positionalReader wraps an io.Reader to that tracks the offset as bytes
// are read. Additionally, it adds a best-effort io.Seeker implementation.
// For a pure io.Reader that is limited to usage of io.SeeekCurrent and
// otherwise fails for seeks relative to the start or end of the stream.
type positionalReader struct {
	reader io.Reader
	pos    int64
}

// Read is a pass-thru to the underlying io.Reader.Read
func (pr *positionalReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.pos += int64(n)
	return
}

// Seek implements io.Seeker. If the wrapped io.Reader also implements
// io.Seeker this is a pass-thru. Otherwise, only io.SeekCurrent is
// supported and ErrUnseekableReader is returned for seeks from start/end.
func (pr *positionalReader) Seek(offset int64, whence int) (int64, error) {
	var err error
	switch s := pr.reader.(type) {
	case io.Seeker:
		pr.pos, err = s.Seek(offset, whence)
	default:
		if whence != io.SeekCurrent {
			err = ErrUnseekableReader
		} else {
			var n int64
			n, err = io.CopyN(ioutil.Discard, pr.reader, offset)
			pr.pos += n
		}
	}
	return pr.pos, err
}

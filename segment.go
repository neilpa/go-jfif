package jfif

import (
	"encoding/binary"
	"io"
)

// SegmentRef represents a distinct region of a JPEG file, identified by a
// 0xff byte followed by a marker byte then optional payload.
type SegmentRef struct {
	// Marker identifies the type of segment.
	Marker
	// Size is the total byte size of a segment, max 2^16 + 2.
	Size int
	// Offset is the address of the 0xff byte at the start of this segment,
	// immediately followed by the marker. Padding and byte stuffing means
	// a preceding segment's size + offset may not equal the next segment's
	// offset.
	Offset int64
}

// Segment provides direct acces to the data bytes of a segment.
type Segment struct {
	SegmentRef
	// Data is the raw bytes of a segment, excluding the initial 4 bytes
	// (e.g. 0xff, marker, and 2-byte length). For segments lacking a
	// length, this will be nil.
	Data []byte
}

// emptyRef creates a segment reference without a payload
func emptyRef(m Marker, off int64) SegmentRef {
	return SegmentRef{m, 2, off}
}

// emptySegment creates a segment without a payload.
func emptySegment(m Marker, off int64) Segment {
	return Segment{ emptyRef(m, off), nil }
}

// decodeRef reads a segment assuming r is positioned at the length of the
// payload after reading the marker.
func decodeRef(r *countReader, m Marker) (SegmentRef, error) {
	// Set the offset to the 0xff byte preceding the marker
	ref := emptyRef(m, r.count - 2)

	// TODO Known empty segments...
	var length uint16
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return ref, err
	}
	if length < 2 {
		return ref, ErrShortSegment
	}

	// Length includes the 2 bytes for itself but not 0xff or marker bytes.
	ref.Size += int(length)
	return ref, nil
}

// decodeSegment reads a segment assuming r is positioned at the length of
// the payload after reading the marker.
func decodeSegment(r *countReader, m Marker) (Segment, error) {
	ref, err := decodeRef(r, m)
	seg := Segment{ ref, nil }
	if err != nil {
		return seg, err
	}
	// Size includes the both marker bytes and 2-byte length
	seg.Data = make([]byte, ref.Size-4)
	_, err = r.Read(seg.Data)
	return seg, err
}

// Load retrieves the payload data for the given segment reference,
func (ref SegmentRef) Load(r io.ReaderAt) (Segment, error) {
	seg := Segment{SegmentRef: ref}
	if ref.Size < 4 {
		return seg, nil
	}

	seg.Data = make([]byte, ref.Size - 4)
	_, err := r.ReadAt(seg.Data, ref.Offset)
	return seg, err
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

// Write serializes the segment to the given writer.
func (s Segment) Write(w io.Writer) error {
	// Everything needs the 0xff, marker and potential payload
	_, err := w.Write([]byte{0xff, byte(s.Marker)})
	if err != nil || s.Data == nil {
		return err
	}
	// Payload size includes it's own 2-bytes
	err = binary.Write(w, binary.BigEndian, uint16(len(s.Data))+2)
	if err != nil {
		return err
	}
	_, err = w.Write(s.Data)
	return err
}

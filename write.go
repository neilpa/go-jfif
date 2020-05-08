package jfif

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"

	xio "neilpa.me/go-x/io"
)

var (
	// ErrUnknownSegment is an attempt to edit an unknown segment
	ErrUnknownSegment = errors.New("Unknown segment")

	// ErrOversizePayload means there's not enough space to update
	// the segment data in-place.
	ErrOversizePayload = errors.New("Oversize payload")

	// ErrOversizeSegment means segment data was to large to append.
	ErrOversizeSegment = errors.New("Oversize segment")
)

// EncodeSegment writes the given segment.
func EncodeSegment(w io.Writer, seg Segment) error { // TODO Segment.WriteTo?
	// Everything else needs the 0xff, marker and potential payload
	_, err := w.Write([]byte{0xff, byte(seg.Marker)})
	if err != nil || seg.Data == nil {
		return err
	}
	// Payload size includes it's own 2-bytes
	// TODO Validate the length of Data here?
	err = binary.Write(w, binary.BigEndian, uint16(len(seg.Data))+2)
	if err != nil {
		return err
	}
	_, err = w.Write(seg.Data)
	return err
}

// Add a new JFIF segment with the given data before the SOS segment. See
// Append for more details.
func Add(path string, m Marker, buf []byte) error {
	return Append(path, Segment{Pointer: Pointer{Marker: m}, Data: buf})
}

// Append new JFIF segments to the file at path. Offsets are ignroed and
// these are added just before SOS segment.
//
// Notes:
//	- Under the hood this creates a temp-copy of the original file so
//		that it can safely insert the new segments in the middle of the
//		file. This avoids potential for corrupting data if an error is
//		hit in the middle of the update. At the end the original path
//		is replaced with a single os.Rename operation.
//
// TODO: Higher-level version of this that could be smarter for XMP data
// TODO: Return the updated pointer data?
func Append(path string, segs ...Segment) error {
	// Prep the buffer for writing
	var buf bytes.Buffer
	for _, seg := range segs {
		l := len(seg.Data) + 2
		if l > math.MaxUint16 {
			return ErrOversizeSegment
		}
		seg.Length = uint16(l) // TODO: what about an embedded Data where the first two bytes are the length

		// TODO Would be nice to avoid yet-another-copy of data and plumb
		// through a custom reader to SpliceFile and the known size.
		if err := EncodeSegment(&buf, seg); err != nil {
			return err
		}
	}

	f, err := os.Open(path)
	ptrs, err := ScanSegments(f)
	if err != nil {
		return err
	}
	last := ptrs[len(ptrs)-1]

	return xio.SpliceFile(f, buf.Bytes(), last.Offset)
}

// File is used to perform in-place updates to JFIF segments to a backing
// file on-disk.
//
// TODO: This may not be all that valuable verse doing a proper splice.
// in a copied version of the file and replacing over top of it. This
// can lead to file corruption if not careful...
type File struct {
	// f is the underlying file on disk.
	f *os.File

	// refs are the minimally scanned segment pointers.
	refs []Pointer
}

// Edit opens and scans segments from a JFIF file. This should be
// used to replace segments in-place without having to re-write the
// full file. Note that this will fail on attempts to write segments
// that would expand beyond the current bounds.
//
// TODO: Otherise, "short-segments" retain the desired size but there
// are 0xFF fill bytes used for padding until the next segment.
func Edit(path string) (*File, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	// TODO: Could track potential padding bytes between segments and
	// use that to "squish" a few segments together without having to
	// re-encode the entire image. Only the parts at the start of it.
	refs, err := ScanSegments(f)
	if err != nil {
		return nil, err
	}
	return &File{f, refs}, nil
}

// Close the underlying file on disk.
func (f *File) Close() (err error) {
	if f.f != nil {
		err = f.f.Close()
		f.f = nil
		f.refs = nil
	}
	return
}

// Query finds existing segments that matches the given marker
func (f *File) Query(m Marker) ([]Pointer, error) {
	refs := make([]Pointer, 0)
	for _, r := range f.refs {
		if r.Marker == m {
			refs = append(refs, r)
		}
	}
	return refs, nil
}

// Update replaces the payload for the given segment ref. Returns an
// error if it's too large or doesn't match a known segment in this
// file.
//
// Note:
//	- This updates the file in-place so all of the general warnings
//		apply w.r.t. potential file corruption. This should be limited
//		to files that have already been copied and are intended to
//		be edited directly.
func (f *File) Update(r Pointer, buf []byte) error {
	var i int
	for ; i < len(f.refs); i++ {
		if f.refs[i] == r {
			break
		}
	}
	if i == len(f.refs) {
		return ErrUnknownSegment
	}

	space := int64(r.Length - 2) // length is inclusive
	if i < len(f.refs)-1 {
		// Potentially more room if there were fill bytes between segments.
		// Account for the 0xFF leader, Marker byte, 2-byte length.
		space = f.refs[i+1].Offset - r.Offset - 4
	}
	if int64(len(buf)) > space || len(buf) > 0xFFFF {
		return ErrOversizePayload
	}

	// Encode the updated segment to disk
	// TODO Need to make sure to update our Pointer copy
	_, err := f.f.Seek(r.Offset, io.SeekStart)
	if err != nil {
		return err
	}

	seg := Segment{ // TODO Can I avoid all the "+/- 2's" everywhere
		Pointer{r.Offset, r.Marker, uint16(len(buf) + 2)},
		buf,
	}
	err = EncodeSegment(f.f, seg)
	if err != nil {
		return err
	}
	if space > int64(len(buf)) {
		return errors.New("TODO: handle the padding case")
	}

	// Update our in-memory location
	f.refs[i] = seg.Pointer
	return nil
}

package jfif

import (
	"errors"
	"io"
	"os"
)

var (
	// ErrUnknownSegment is an attempt to edit an unknown segment
	ErrUnknownSegment = errors.New("Unknown segment")

	// ErrOversizePayload means there's not enough space to update
	// the segment data in-place.
	ErrOversizePayload = errors.New("Oversize payload")
)

// File is used to perform in-place updates to JFIF segments to a backing
// file on-disk.
type File struct {
	// f is the underlying file on disk.
	f *os.File

	// refs are the intially scanned segment pointers.
	refs []SegmentP
}

// Edit opens and scans segments from a JFIF file. This should be
// used to replace segments in-place without having to re-write the
// full file. Note that this will fail on attempts to write segments
// that would expend beyond the current bounds. Otherise, "short-segments"
// retain the desired size but there are 0xFF fill bytes used for padding
// until the next segment.
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
func (f *File) Query(m Marker) ([]SegmentP, error) {
	refs := make([]SegmentP ,0)
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
func (f *File) Update(r SegmentP, buf []byte) error {
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
	// TODO Need to make sure to update our SegmentP copy
	_, err := f.f.Seek(r.Offset, io.SeekStart)
	if err != nil {
		return err
	}

	seg := Segment{ // TODO Can I avoid all the "+/- 2's" everywhere
		SegmentP{r.Offset, r.Marker, uint16(len(buf)+2)}, buf,
	}
	err = EncodeSegment(f.f, seg)
	if err != nil {
		return err
	}
	if space > int64(len(buf)) {
		return errors.New("TODO: handle the padding case")
	}

	// Update our in-memory location
	f.refs[i] = seg.SegmentP
	return nil
}

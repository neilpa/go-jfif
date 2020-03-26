package jfif

import "fmt"

// Marker identifies the various types of JFIF segments.
type Marker byte

// https://www.disktuna.com/list-of-jpeg-markers/
const (
	// SOF0 (Start Of Frame 0) marker indicates a baseline DCT.
	SOF0 Marker = iota + 0xC0
	// SOF1 (Start Of Frame 1) marker indicates an extended sequential DCT.
	SOF1
	// SOF2 (Start Of Frame 2) marker indicates a progressive DCT.
	SOF2
	// SOF3 (Start Of Frame 3) marker indicates a lossless (sequential) image.
	SOF3
	// DHT (Define Huffman Table[s]) marker specifies one or more Huffman tables.
	DHT
	// SOF5 (Start Of Frame 5) marker indicates a differential sequential DCT.
	SOF5
	// SOF6 (Start Of Frame 6) marker indicates a differential progressive DCT.
	SOF6
	// SOF7 (Start Of Frame 7) marker indicates a differential lossless (sequential) image.
	SOF7
	// JPG (JPEG Extensions) marker.
	JPG
	// SOF9 (Start Of Frame 9) marker indicates extended sequentital DCT with arithmetic coding.
	SOF9
	// SOF10 (Start Of Frame 10) marker indicates progressive DCT with arithmetic coding.
	SOF10
	// SOF11 (Start Of Frame 11) marker indicates lossless (sequential) with arithmetic coding.
	SOF11
	// DAC (Define Arithmetic Coding) marker specifies one or more arithmetic coding tables.
	DAC
	// SOF13 (Start Of Frame 13) marker indicates a differential sequential DCT with arithmetic coding.
	SOF13
	// SOF14 (Start Of Frame 14) marker indicates a differential progressive DCT with arithmetic coding.
	SOF14
	// SOF15 (Start Of Frame 15) marker indicates a differential lossless (sequential) DCT with arithmetic coding.
	SOF15

	// RST0 (Restart 0) marker.
	RST0
	// RST1 (Restart 1) marker.
	RST1
	// RST2 (Restart 2) marker.
	RST2
	// RST3 (Restart 3) marker.
	RST3
	// RST4 (Restart 4) marker.
	RST4
	// RST5 (Restart 5) marker.
	RST5
	// RST6 (Restart 6) marker.
	RST6
	// RST7 (Restart 7) marker.
	RST7

	// SOI (Start Of Image) marker
	SOI
	// EOI (End Of Image) marker.
	EOI
	// SOS (Start Of Scan/Stream) marker begins a top-to-bottom scan of the image. In baseline DCT JPEG images, there is generally a single scan. Progressive DCT JPEG images usually contain multiple scans. This marker specifies which slice of data it will contain, and is immediately followed by entropy-coded data.
	SOS

	// DQT (Define Quantization Table[s]) marker specifies one or more quantization tables.
	DQT
	// DNL (Define Numer of Lines) marker (uncommon).
	DNL
	// DRI (Define Restart Interval) marker specifies the interval between RSTn markers, in Minimum Coded Units (MCUs). This marker is followed by two bytes indicating the fixed size so it can be treated like any other variable size segment.
	DRI
	// DHP (Define Hierarchical Progression) marker (uncommon).
	DHP
	// EXP (Expand Reference Component) marker (uncommon).
	EXP

	// APP0 (Application-specific 0) marker for custom metadata (JFIF; AVI1).
	APP0
	// APP1 (Application-specific 1) marker for custom metadata (Exif; TIFF IFD; JPEG Thumbnail; Adobe XMP).
	APP1
	// APP2 (Application-specific 2) marker for custom metadata (ICC color profile; FlashPix).
	APP2
	// APP3 (Application-specific 3) marker for custom metadata (uncommon; JPS Tag for Sterioscopic JPEG).
	APP3
	// APP4 (Application-specific 4) marker for custom metadata (uncommon).
	APP4
	// APP5 (Application-specific 5) marker for custom metadata (uncommon).
	APP5
	// APP6 (Application-specific 6) marker for custom metadata (uncommon; NITF Lossless profile).
	APP6
	// APP7 (Application-specific 7) marker for custom metadata (uncommon).
	APP7
	// APP8 (Application-specific 8) marker for custom metadata (uncommon).
	APP8
	// APP9 (Application-specific 9) marker for custom metadata (uncommon).
	APP9
	// APP10 (Application-specific 10) marker for custom metadata (uncommon; ActiveObject).
	APP10
	// APP11 (Application-specific 11) marker for custom metadata (uncommon; HELIOS JPEG Resources).
	APP11
	// APP12 (Application-specific 12) marker for custom metadata (Old digicams picture info, Photoshop Save for Web: Ducky).
	APP12
	// APP13 (Application-specific 13) marker for custom metadata (Photoshop Save As: IRB, 8BIM, IPTC).
	APP13
	// APP14 (Application-specific 14) marker for custom metadata (uncommon).
	APP14
	// APP15 (Application-specific 15) marker for custom metadata (uncommon).
	APP15

	// JPG0 (JPEG Extension 0) marker (uncommon).
	JPG0
	// JPG1 (JPEG Extension 1) marker (uncommon).
	JPG1
	// JPG2 (JPEG Extension 2) marker (uncommon).
	JPG2
	// JPG3 (JPEG Extension 3) marker (uncommon).
	JPG3
	// JPG4 (JPEG Extension 4) marker (uncommon).
	JPG4
	// JPG5 (JPEG Extension 5) marker (uncommon).
	JPG5
	// JPG6 (JPEG Extension 6) marker (uncommon).
	JPG6
	// JPG7 (JPEG Extension 7) marker (Lossless JPEG).
	JPG7
	// JPG8 (JPEG Extension 8) marker (Lossless JPEG Extension Parameters).
	JPG8
	// JPG9 (JPEG Extension 9) marker (uncommon).
	JPG9
	// JPG10 (JPEG Extension 10) marker (uncommon).
	JPG10
	// JPG11 (JPEG Extension 11) marker (uncommon).
	JPG11
	// JPG12 (JPEG Extension 12) marker (uncommon).
	JPG12
	// JPG13 (JPEG Extension 13) marker (uncommon).
	JPG13

	// COM (Comment) marker contains a text comment.
	COM

	// SOF48 is the JPEG-LS marker.
	SOF48 = JPG7
	// LSE is the JPEG-LS Extension marker.
	LSE = JPG8
)

var markerTags = []string{
	"SOF0",
	"SOF1",
	"SOF2",
	"SOF3",
	"DHT",
	"SOF5",
	"SOF6",
	"SOF7",
	"JPG",
	"SOF9",
	"SOF10",
	"SOF11",
	"DAC",
	"SOF13",
	"SOF14",
	"SOF15",
	"RST0",
	"RST1",
	"RST2",
	"RST3",
	"RST4",
	"RST5",
	"RST6",
	"RST7",
	"SOI",
	"EOI",
	"SOS",
	"DQT",
	"DNL",
	"DRI",
	"DHP",
	"EXP",
	"APP0",
	"APP1",
	"APP2",
	"APP3",
	"APP4",
	"APP5",
	"APP6",
	"APP7",
	"APP8",
	"APP9",
	"APP10",
	"APP11",
	"APP12",
	"APP13",
	"APP14",
	"APP15",
	"JPG0",
	"JPG1",
	"JPG2",
	"JPG3",
	"JPG4",
	"JPG5",
	"JPG6",
	"SOF48",
	"LSE",
	"JPG9",
	"JPG10",
	"JPG11",
	"JPG12",
	"JPG13",
	"COM",
}

func (m Marker) String() string {
	i := int(m - SOF0)
	if 0 <= i && i < len(markerTags) {
		return markerTags[i]
	}
	return fmt.Sprintf("0x%X", byte(m))
}

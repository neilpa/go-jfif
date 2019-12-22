package jfif

import "fmt"

// Marker identifies the various types of JFIF segments.
type Marker byte

// TODO https://www.disktuna.com/list-of-jpeg-markers/
const (
	// SOI (Start Of Image) marker
	SOI = 0xD8

	// SOF0 (Start Of Frame 0) marker indicates a baseline DCT.
	SOF0 = 0xC0
	// SOF2 (Start Of Frame 2) marker indicates a progressive DCT.
	SOF2 = 0xC2

	// DHT (Define Huffman Table[s]) marker specifies one or more Huffman tables.
	DHT = 0xC4
	// DQT (Define Quantization Table[s]) marker specifies one or more quantization tables.
	DQT = 0xDB

	// DRI (Define Restart Interval) marker specifies the interval between RSTn markers, in Minimum Coded Units (MCUs). This marker is followed by two bytes indicating the fixed size so it can be treated like any other variable size segment.
	DRI = 0xDD

	// SOS (Start Of Scan/Stream) marker begins a top-to-bottom scan of the image. In baseline DCT JPEG images, there is generally a single scan. Progressive DCT JPEG images usually contain multiple scans. This marker specifies which slice of data it will contain, and is immediately followed by entropy-coded data.
	SOS = 0xDA

	// RST0 (Restart 0) marker.
	RST0 = 0xD0
	// RST1 (Restart 1) marker.
	RST1 = 0xD1
	// RST2 (Restart 2) marker.
	RST2 = 0xD2
	// RST3 (Restart 3) marker.
	RST3 = 0xD3
	// RST4 (Restart 4) marker.
	RST4 = 0xD4
	// RST5 (Restart 5) marker.
	RST5 = 0xD5
	// RST6 (Restart 6) marker.
	RST6 = 0xD6
	// RST7 (Restart 7) marker.
	RST7 = 0xD7

	// APP0 (Application-specific 0) marker for custom metadata.
	APP0 = 0xE0
	// APP1 (Application-specific 1) marker for custom metadata.
	APP1 = 0xE1
	// APP2 (Application-specific 2) marker for custom metadata.
	APP2 = 0xE2
	// APP3 (Application-specific 3) marker for custom metadata.
	APP3 = 0xE3
	// APP4 (Application-specific 4) marker for custom metadata.
	APP4 = 0xE4
	// APP5 (Application-specific 5) marker for custom metadata.
	APP5 = 0xE5
	// APP6 (Application-specific 6) marker for custom metadata.
	APP6 = 0xE6
	// APP7 (Application-specific 7) marker for custom metadata.
	APP7 = 0xE7
	// APP8 (Application-specific 8) marker for custom metadata.
	APP8 = 0xE8
	// APP9 (Application-specific 9) marker for custom metadata.
	APP9 = 0xE9
	// APP10 (Application-specific 10) marker for custom metadata.
	APP10 = 0xEA
	// APP11 (Application-specific 11) marker for custom metadata.
	APP11 = 0xEB
	// APP12 (Application-specific 12) marker for custom metadata.
	APP12 = 0xEC
	// APP13 (Application-specific 13) marker for custom metadata.
	APP13 = 0xED
	// APP14 (Application-specific 14) marker for custom metadata.
	APP14 = 0xEE
	// APP15 (Application-specific ) marker for custom metadata.
	APP15 = 0xEF

	// COM (Comment) marker contains a text comment.
	COM = 0xFE

	// EOI (End Of Image) marker.
	EOI = 0xD9
)

func (m Marker) String() string {
	switch m {
	case SOI:
		return "SOI"
	case SOF0:
		return "SOF0"
	case SOF2:
		return "SOF2"
	case DHT:
		return "DHT"
	case DQT:
		return "DQT"
	case DRI:
		return "DRI"
	case SOS:
		return "SOS"
	case RST0:
		return "RST0"
	case RST1:
		return "RST1"
	case RST2:
		return "RST2"
	case RST3:
		return "RST3"
	case RST4:
		return "RST4"
	case RST5:
		return "RST5"
	case RST6:
		return "RST6"
	case RST7:
		return "RST7"
	case APP0:
		return "APP0"
	case APP1:
		return "APP1"
	case APP2:
		return "APP2"
	case APP3:
		return "APP3"
	case APP4:
		return "APP4"
	case APP5:
		return "APP5"
	case APP6:
		return "APP6"
	case APP7:
		return "APP7"
	case APP8:
		return "APP8"
	case APP9:
		return "APP9"
	case APP10:
		return "APP10"
	case APP11:
		return "APP11"
	case APP12:
		return "APP12"
	case APP13:
		return "APP13"
	case APP14:
		return "APP14"
	case APP15:
		return "APP15"
	case COM:
		return "COM"
	case EOI:
		return "EOI"
	}
	return fmt.Sprintf("%X", byte(m))
}

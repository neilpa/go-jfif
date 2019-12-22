package jfif

// Various signatures that identify different APPn segments

// APP0
const (
	SigJFIF = "JFIF\x00"
	SigJFXX = "JFXX\x00"
)

// APP1
const (
	SigExif        = "Exif\x00\x00"
	SigXMP         = "http://ns.adobe.com/xap/1.0/\x00"
	SigExtendedXMP = "http://ns.adobe.com/xmp/extension/\x00"
)

// APP13
const (
	SigPhotoshop3 = "Photoshop 3.0\x00"
	SigPhotoshop2 = "Adobe_Photoshop2.5:"
)

var appnSigs = [16][]string{
	{SigJFIF, SigJFXX},
	{SigExif, SigXMP, SigExtendedXMP},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{SigPhotoshop3, SigPhotoshop2},
}

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
	SigExifAlt     = "Exif\x00\xFF"
	SigXMP         = "http://ns.adobe.com/xap/1.0/\x00"
	SigExtendedXMP = "http://ns.adobe.com/xmp/extension/\x00"
)

// APP2
const (
	SigICCProfile = "ICC_PROFILE\x00"
)

// APP3
const (
	SigMETA = "META\x00\x00"
	SigMeta = "Meta\x00\x00"
)

// APP12
const (
	SigDucky = "Ducky\x00"
)

// APP13
const (
	SigPhotoshop3 = "Photoshop 3.0\x00"
	SigPhotoshop2 = "Adobe_Photoshop2.5:"
)

// APP14
const (
	SigAdobe = "Adobe\x00"
)


var appnSigs = [16][]string{
	{SigJFIF, SigJFXX},
	{SigExif, SigXMP, SigExtendedXMP},
	{SigICCProfile},
	{SigMETA, SigMeta},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{},
	{SigDucky},
	{SigPhotoshop3, SigPhotoshop2},
	{SigAdobe},
}

# go-jfif

[![CI](https://github.com/neilpa/go-jfif/workflows/CI/badge.svg)](https://github.com/neilpa/go-jfif/actions/)

Basic reading of [JPEG File Interchange Format (JFIF)][wiki-jfif] segments

[wiki-jfif]: https://en.wikipedia.org/wiki/JPEG_File_Interchange_Format#File_format_structure

## Usage

TODO

## Tools

There are a few simple CLI tools include for manipulating JPEG/JFIF files that exercise functionality
exposed by this module.

### jfifcom

Append free-form comment (`COM`) segments to an existing JPEG file.

### jfifstat

Prints segment markers, sizes, and optional `APPN` signatures from a JPEG until the image stream.

### xmpdump

Extracts and prints [XMP](https://www.adobe.com/products/xmp.html) data from `APP1` segments from JPEG file(s).

## License

[MIT](/LICENSE)

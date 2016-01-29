// Package ioutil implements some I/O utility functions that
// are UTF-encoding agnostic.
package utfutil

// These functions autodetect UTF BOM and return UTF-8. If no
// BOM is found, they assume UTF-8 (Linux) or UTF-16LE (Windows).
// You can use them as replacements for os.Open() and ioutil.ReadFile()
// when the encoding of the file is unknown.

// Since it is impossible to guess 100% correctly if there is no BOM,
// the functions take a 2nd parameter of type "EncodingHint" where you
// specify the default encoding.

// In the future I'd like to add a default type "AUTO" which
// makes some educated guesses about the default type, similar
// to the uchardet command. Hopefully that kind of functionality
// will be added to golang.org/x/text/encoding/unicode :-)

// Inspiration: I wrote this after spending half a day trying
// to figure out how to use unicode.BOMOverride.
// Hopefully this will save other golang newbies from the same.

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// EncodingHint indicates the file's encoding if there is no BOM.
type EncodingHint int

const (
	UTF8 EncodingHint = iota
	UTF16LE
	UTF16BE
	WINDOWS = UTF16LE // File came from a MS-Windows system
	POSIX   = UTF8    // File came from Unix or Unix-like systems
	HTML5   = UTF8    // File came from the web
)

// About utfutil.HTML5:
// This technique is recommended by the W3C for use in HTML 5:
// "For compatibility with deployed content, the byte order
// mark (also known as BOM) is considered more authoritative
// than anything else." http://www.w3.org/TR/encoding/#specification-hooks

// NewReader wraps a Reader to decode Unicode to UTF-8 as it reads.
func NewReader(r io.Reader, d EncodingHint) io.Reader {
	var decoder *encoding.Decoder
	switch d {
	case UTF8:
		// Make a transformer that assumes UTF-8 but abides by the BOM.
		decoder = unicode.UTF8.NewDecoder()
	case UTF16LE:
		// Make an tranformer that decodes MS-Windows (16LE) UTF files:
		winutf := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		// Make a transformer that is like winutf, but abides by BOM if found:
		decoder = winutf.NewDecoder()
	case UTF16BE:
		// Make an tranformer that decodes UTF-16BE files:
		utf16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		// Make a transformer that is like utf16be, but abides by BOM if found:
		decoder = utf16be.NewDecoder()
	}

	// Make a Reader that uses utf16bom:
	return transform.NewReader(r, unicode.BOMOverride(decoder))
}

// OpenFile is the equivalent of os.Open().
func OpenFile(name string, d EncodingHint) (io.Reader, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return NewReader(f, d), nil
}

// ReadFile is the equivalent of ioutil.ReadFile()
func ReadFile(name string, d EncodingHint) ([]byte, error) {
	file, err := OpenFile(name, d)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

// BytesReader is a convenience function that takes a []byte and decodes them to UTF-8.
func BytesReader(b []byte, d EncodingHint) io.Reader {
	return NewReader(bytes.NewReader(b), d)
}

// NewScanner is a convenience function that takes a filename and returns a scanner.
func NewScanner(name string, d EncodingHint) (*bufio.Scanner, error) {
	f, err := OpenFile(name, d)
	if err != nil {
		return nil, err
	}
	return bufio.NewScanner(f), nil
}

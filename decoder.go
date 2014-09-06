// Package ini provides parse capabilities for .ini files.
package ini

import (
	"github.com/mitchellh/mapstructure"
	"io"
	"os"
)

// Alias for map[string]map[string]string
type Config config

// Struct to parse .ini format from an io.reader
type Decoder struct {
	rd      io.Reader
	options Options
}

// Struct to contain options for ini.Decoder
type Options struct {
	IdRegexp     string
	SepChars     []byte
	CommentChars []byte
	LowCaseIds   bool
}

// Default options for ini.Decoder
var DefaultOptions Options = Options{
	IdRegexp:     idDefaultRegex,
	SepChars:     []byte{'='},
	CommentChars: []byte{';'},
	LowCaseIds:   true,
}

// Creates a new ini.Decoder from an io.Reader
func NewDecoder(rd io.Reader) *Decoder {
	return &Decoder{rd, DefaultOptions}
}

// Creates a new ini.Decoder from an io.Reader with custom options
func NewDecoderWithOptions(rd io.Reader, opts Options) *Decoder {
	return &Decoder{rd, opts}
}

// Set the separator characters between keys and values. Defaults to '='
func (d *Decoder) SepChars(sepChars []byte) {
	d.options.SepChars = sepChars
}

// Set the character used to start a comment. Defaults to ';'
func (d *Decoder) CommentChars(commentChars []byte) {
	d.options.CommentChars = commentChars
}

// Set if the keys should be converted to lower case. Defaults to true.
func (d *Decoder) LowCaseIds(lowCaseIds bool) {
	d.options.LowCaseIds = lowCaseIds
}

// Set the regexp to check if the key is valid. Defaults to: "^[a-z][a-z0-9_]+$"
func (d *Decoder) IdRegexp(idRegexp string) {
	d.options.IdRegexp = idRegexp
}

// Decode the io.Reader contained into the given interface.
// Returns an error on failure
func (d *Decoder) Decode(r interface{}) error {
	pars := newParserWithOptions(d.rd,
		d.options.IdRegexp, d.options.LowCaseIds,
		d.options.SepChars, d.options.CommentChars)
	if err := pars.parseConfig(); err != nil {
		return err
	}
	if err := mapstructure.WeakDecode(pars.currentConfig, r); err != nil {
		return err
	}
	return nil
}

// Decode the given file to the given interface
func DecodeFile(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return NewDecoder(file).Decode(v)
}

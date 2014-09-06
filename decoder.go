package ini

import (
	"github.com/mitchellh/mapstructure"
	"io"
	"os"
)

type Config config

type Decoder struct {
	rd      io.Reader
	options Options
}

type Options struct {
	IdRegexp     string
	SepChars     []byte
	CommentChars []byte
	LowCaseIds   bool
}

var DefaultOptions Options = Options{
	IdRegexp:     idDefaultRegex,
	SepChars:     []byte{'='},
	CommentChars: []byte{';'},
	LowCaseIds:   true,
}

func NewDecoder(rd io.Reader) *Decoder {
	return &Decoder{rd, DefaultOptions}
}

func NewDecoderWithOptions(rd io.Reader, opts Options) *Decoder {
	return &Decoder{rd, opts}
}

func (d *Decoder) SepChars(sepChars []byte) {
	d.options.SepChars = sepChars
}

func (d *Decoder) CommentChars(commentChars []byte) {
	d.options.CommentChars = commentChars
}

func (d *Decoder) LowCaseIds(lowCaseIds bool) {
	d.options.LowCaseIds = lowCaseIds
}

func (d *Decoder) IdRegexp(idRegexp string) {
	d.options.IdRegexp = idRegexp
}

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

func DecodeFile(path string, v interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return NewDecoder(file).Decode(v)
}

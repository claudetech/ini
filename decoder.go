package ini

import (
	"github.com/mitchellh/mapstructure"
	"io"
)

type Decoder struct {
	rd      io.Reader
	options Options
}

type Options struct {
	IdRegexp     string
	SepChars     []byte
	commentChars []byte
	LowCaseIds   bool
}

var defaultOptions Options = Options{
	IdRegexp:     idDefaultRegex,
	SepChars:     []byte{'='},
	commentChars: []byte{';'},
	LowCaseIds:   true,
}

func NewDecoder(rd io.Reader) *Decoder {
	return &Decoder{rd, defaultOptions}
}

func NewDecoderWithOptions(rd io.Reader, opts Options) *Decoder {
	return &Decoder{rd, opts}
}

func (d *Decoder) SepChars(sepChars []byte) {
	d.options.SepChars = sepChars
}

func (d *Decoder) CommentChars(commentChars []byte) {
	d.options.commentChars = commentChars
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
		d.options.SepChars, d.options.commentChars)
	if err := pars.parseConfig(); err != nil {
		return err
	}
	if err := mapstructure.Decode(r, pars.currentConfig); err != nil {
		return err
	}
	return nil
}

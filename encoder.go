package ini

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
)

type Encoder struct {
	w       *bufio.Writer
	sepChar byte
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: bufio.NewWriter(w), sepChar: '='}
}

func (e *Encoder) SepChar(c byte) {
	e.sepChar = c
}

func (e *Encoder) writeSection(section string, conf map[string]string) error {
	s := fmt.Sprintf("[%s]\n", section)
	if _, err := e.w.WriteString(s); err != nil {
		return err
	}
	for key, val := range conf {
		entry := fmt.Sprintf("%s %c %s\n", key, e.sepChar, val)
		if _, err := e.w.WriteString(entry); err != nil {
			return err
		}
	}
	return e.w.Flush()
}

func (e *Encoder) Encode(v interface{}) error {
	var conf Config
	if err := mapstructure.Decode(v, &conf); err != nil {
		return err
	}
	for section, values := range conf {
		e.writeSection(section, values)
	}
	return nil
}

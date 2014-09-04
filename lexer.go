package ini

import (
	"bufio"
	"bytes"
	"io"
)

type lexer struct {
	reader       *bufio.Reader
	sepChars     []byte
	commentChars []byte
}

func newLexer(reader io.Reader) *lexer {
	return newLexerWithOptions(bufio.NewReader(reader), []byte{'='}, []byte{';'})
}

func newLexerWithOptions(reader io.Reader, sepChars []byte, commentChars []byte) *lexer {
	return &lexer{bufio.NewReader(reader), sepChars, commentChars}
}

func (l *lexer) peekNext() (byte, error) {
	bytes, err := l.reader.Peek(1)
	if err != nil {
		return '0', nil
	}
	return bytes[0], nil
}

func (l *lexer) nextToken() (token, error) {
	nextByte, err := l.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case nextByte == ' ' || nextByte == '\t':
		return &spaceToken{string(nextByte)}, nil
	case bytes.IndexByte(l.sepChars, nextByte) > -1:
		return &sepToken{}, nil
	case bytes.IndexByte(l.commentChars, nextByte) > -1:
		return &commentToken{}, nil
	case nextByte == '\n' || nextByte == '\r':
		if nextByte == '\r' {
			if n, err := l.peekNext(); err == nil && n == '\n' {
				_, _ = l.reader.ReadByte()
			}
		}
		return &newLineToken{}, nil
	case nextByte == '[' || nextByte == ']' || nextByte == '"':
		return &symbolToken{string(nextByte)}, nil
	default:
		return &otherToken{string(nextByte)}, nil
	}
}

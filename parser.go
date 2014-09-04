package ini

import (
	"io"
)

type parser struct {
	lex          *lexer
	currentToken token
}

func newParser(rd io.Reader) *parser {
	return &parser{lex: newLexer(rd)}
}

func newParserWithOptions(rd io.Reader, sepChars []byte, commentChars []byte) *parser {
	return &parser{lex: newLexerWithOptions(rd, sepChars, commentChars)}
}

func (p *parser) eat(token) {

}

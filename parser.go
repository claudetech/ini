package ini

import (
	"fmt"
	"io"
)

type parseError struct {
	p       *parser
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("Parse error at %d:%d. %s",
		e.p.currentLine, e.p.currentChar)
}

type parser struct {
	lex          *lexer
	currentToken token
	currentLine  int
	currentChar  int
}

func newParser(rd io.Reader) *parser {
	return &parser{lex: newLexer(rd), currentLine: 1, currentChar: 0}
}

func newParserWithOptions(rd io.Reader, sepChars []byte, commentChars []byte) *parser {
	return &parser{lex: newLexerWithOptions(rd, sepChars, commentChars)}
}

func (p *parser) eat(typ tokenType) (t token, err error) {
	if t = p.advance(); t == nil {
		return nil, nil
	}
	if typ != p.currentToken.getType() {
		msg := fmt.Sprintf("Expected %s, got %s.", typ.ToString(), p.currentToken.getType().ToString())
		err = parseError{p, msg}
	}
	return
}

func (p *parser) advance() token {
	tok, err := p.lex.nextToken()
	if err != nil {
		// EOF
		return nil
	}
	if tok.getType() == newLineTokType {
		p.currentChar = 0
		p.currentLine += 1
	} else {
		p.currentChar += 1
	}

	p.currentToken = tok
	return p.currentToken
}

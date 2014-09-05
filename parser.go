package ini

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	idDefaultRegex = "[a-z][a-z0-9_]+"
)

type parseError struct {
	p       *parser
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("Parse error at %d:%d. %s",
		e.p.currentLine, e.p.currentChar, e.message)
}

func newTokenError(p *parser, expected string) parseError {
	var value string
	tok := p.currentToken
	switch t := tok.(type) {
	case *symbolToken:
		value = t.symbol
	case *otherToken:
		value = t.value
	default:
		value = ""
	}

	dispVal := " "
	if value != "" {
		dispVal += "'" + value + "'"
	}

	msg := fmt.Sprintf("Expected %s, got %s%s.", expected, tok.getType().ToString(), dispVal)
	return parseError{p, msg}
}

type parser struct {
	lex          *lexer
	currentToken token
	currentLine  int
	currentChar  int
	idRegexp     *regexp.Regexp
	lowCaseIds   bool
}

func makeParser(lex *lexer, regex string, lowCaseIds bool) *parser {
	idRegexp, err := regexp.Compile(regex)
	if err != nil {
		idRegexp, _ = regexp.Compile(idDefaultRegex)
	}
	return &parser{
		lex:          lex,
		currentToken: nil,
		currentLine:  1,
		currentChar:  0,
		idRegexp:     idRegexp,
		lowCaseIds:   lowCaseIds,
	}
}

func newParser(rd io.Reader) *parser {
	lex := newLexer(rd)
	return makeParser(lex, idDefaultRegex, true)
}

func newParserWithOptions(rd io.Reader,
	idRegexp string, lowCaseIds bool,
	sepChars []byte, commentChars []byte) *parser {
	lex := newLexerWithOptions(rd, sepChars, commentChars)
	return makeParser(lex, idRegexp, lowCaseIds)
}

func (p *parser) eat(typ tokenType) (t token, err error) {
	t = p.currentToken
	if t != nil && typ != p.currentToken.getType() {
		err = newTokenError(p, typ.ToString())
	}
	p.advance()
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

func (p *parser) parseSection() (sectionName string, err error) {
	if s, err := p.eat(symbolTokType); err != nil || s.(*symbolToken).symbol != "[" {
		return "", newTokenError(p, "[")
	}

	token := p.currentToken

	for t, ok := token.(*otherToken); ok; t, ok = token.(*otherToken) {
		v := t.value
		if p.lowCaseIds {
			v = strings.ToLower(v)
		}
		sectionName += v
		token = p.advance()
	}

	if s, err := p.eat(symbolTokType); err != nil || s.(*symbolToken).symbol != "]" {
		return "", newTokenError(p, "]")
	}
	if !p.idRegexp.MatchString(sectionName) {
		msg := fmt.Sprintf("Bad section name: %s. Should match %s.",
			sectionName, p.idRegexp.String())
		return "", parseError{p, msg}
	}
	return sectionName, nil
}

package ini

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type config map[string]map[string]string

const (
	idDefaultRegex = "^[a-z][a-z0-9_]+$"
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
	tok := p.currentToken
	value := stringValue(tok)

	dispVal := " "
	if value != "" {
		dispVal += "'" + value + "'"
	}

	msg := fmt.Sprintf("Expected %s, got %s%s.", expected, tok.getType().String(), dispVal)
	return parseError{p, msg}
}

type parser struct {
	lex            *lexer
	currentToken   token
	currentLine    int
	currentChar    int
	idRegexp       *regexp.Regexp
	lowCaseIds     bool
	currentSection string
	currentConfig  config
}

func makeParser(lex *lexer, regex string, lowCaseIds bool) *parser {
	idRegexp, err := regexp.Compile(regex)
	if err != nil {
		idRegexp, _ = regexp.Compile(idDefaultRegex)
	}
	parser := &parser{
		lex:            lex,
		currentToken:   nil,
		currentLine:    1,
		currentChar:    0,
		currentSection: "",
		currentConfig:  make(map[string]map[string]string),
		idRegexp:       idRegexp,
		lowCaseIds:     lowCaseIds,
	}
	parser.advance()
	return parser
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
		err = newTokenError(p, typ.String())
	}
	p.advance()
	return
}

func (p *parser) advance() token {
	tok, err := p.lex.nextToken()
	if err != nil {
		// EOF
		p.currentToken = nil
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

func (p *parser) parseIdentifier() (ident string, err error) {
	var buffer bytes.Buffer

	shouldStop := func(tokType tokenType) bool {
		return tokType == commentTokType || tokType == symbolTokType || tokType == sepTokType
	}

	for token := p.currentToken; token != nil && !shouldStop(token.getType()); token = p.advance() {
		v := stringValue(token)
		if p.lowCaseIds {
			v = strings.ToLower(v)
		}
		buffer.WriteString(v)
	}
	ident = strings.TrimRight(buffer.String(), " \t")

	if !p.idRegexp.MatchString(ident) {
		msg := fmt.Sprintf("Bad key name: %s. Should match %s.",
			ident, p.idRegexp.String())
		err = parseError{p, msg}
	}
	return
}

func (p *parser) parseSection() (sectionName string, err error) {
	if s, err := p.eat(symbolTokType); err != nil || s.(*symbolToken).symbol != "[" {
		return "", newTokenError(p, "[")
	}

	if sectionName, err = p.parseIdentifier(); err != nil {
		return
	}

	if s, err := p.eat(symbolTokType); err != nil || s.(*symbolToken).symbol != "]" {
		err = newTokenError(p, "]")
	}

	return
}

func (p *parser) parseValue() (value string, err error) {
	var buffer bytes.Buffer
	token := p.currentToken
	for token != nil && token.getType() != newLineTokType && token.getType() != commentTokType {
		buffer.WriteString(stringValue(token))
		token = p.advance()
	}
	value = strings.TrimRight(buffer.String(), " \t")
	return
}

func (p *parser) skipSpaces() {
	for token := p.currentToken; token != nil && token.getType() == spaceTokType; token = p.advance() {
	}
}

func (p *parser) parseAssignment() (ident string, value string, err error) {
	ident, err = p.parseIdentifier()
	if err != nil {
		return
	}
	p.skipSpaces()
	if _, err = p.eat(sepTokType); err != nil {
		return
	}
	p.skipSpaces()
	value, err = p.parseValue()
	if err != nil {
		return
	}
	return
}

func (p *parser) skipComment() {
	for token := p.currentToken; token != nil && token.getType() != newLineTokType; token = p.advance() {
	}
}

func (p *parser) changeSection() error {
	sec, err := p.parseSection()
	if err != nil {
		return err
	}
	p.currentSection = sec
	if _, ok := p.currentConfig[p.currentSection]; !ok {
		p.currentConfig[p.currentSection] = make(map[string]string)
	}
	return nil
}

func (p *parser) makeAssignement() error {
	key, value, err := p.parseAssignment()
	if err != nil {
		return err
	}
	p.currentConfig[p.currentSection][key] = value
	return nil
}

func (p *parser) parseLine() (err error) {
	p.skipSpaces()
	if p.currentToken == nil {
		return nil
	}
	switch t := p.currentToken.(type) {
	case *symbolToken:
		if t.symbol != "[" {
			return parseError{p, fmt.Sprintf("Unexpected token %s.", t.symbol)}
		}
		err = p.changeSection()
	case *otherToken:
		if p.currentSection == "" {
			return parseError{p, "Expected section start"}
		}
		err = p.makeAssignement()
	case *commentToken:
		p.skipComment()
	case *sepToken:
		return parseError{p, "Unexpected separator token."}
	default:
	}
	if err != nil {
		return err
	}
	p.skipSpaces()
	if p.currentToken != nil && p.currentToken.getType() == commentTokType {
		p.skipComment()
	}
	_, err = p.eat(newLineTokType)
	return
}

func (p *parser) parseConfig() (err error) {
	for err = p.parseLine(); p.currentToken != nil; err = p.parseLine() {
		if err != nil {
			return
		}
	}
	return
}

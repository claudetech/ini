package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strings"
)

func getToken(lex *lexer) token {
	token, err := lex.nextToken()
	Expect(err).To(BeNil())
	return token
}

var _ = Describe("lexer", func() {
	Describe("NextToken", func() {
		It("should return spaces and tabs", func() {
			lex := newLexer(strings.NewReader(" \t "))
			expected := []string{" ", "\t", " "}
			for _, e := range expected {
				token := getToken(lex)
				Expect(token.(*spaceToken).value).To(Equal(e))
			}
		})

		It("should return new lines", func() {
			lex := newLexer(strings.NewReader("\n\r\n\r"))
			for i := 0; i < 3; i++ {
				token := getToken(lex)
				Expect(func() { _ = token.(*newLineToken) }).NotTo(Panic())
			}
		})

		It("should return symbols", func() {
			lex := newLexer(strings.NewReader("[]\""))
			expected := []string{"[", "]", "\""}
			for _, e := range expected {
				token := getToken(lex)
				Expect(token.(*symbolToken).value).To(Equal(e))
			}
		})

		It("should return seps", func() {
			lex := newLexer(strings.NewReader("="))
			Expect(func() { _ = getToken(lex).(*sepToken) }).NotTo(Panic())
			lex = newLexerWithOptions(strings.NewReader(":"), []byte{':'}, []byte{';'})
			Expect(func() { _ = getToken(lex).(*sepToken) }).NotTo(Panic())
		})

		It("should return comments", func() {
			lex := newLexer(strings.NewReader(";"))
			Expect(func() { _ = getToken(lex).(*commentToken) }).NotTo(Panic())
			lex = newLexerWithOptions(strings.NewReader("#"), []byte{'='}, []byte{'#'})
			Expect(func() { _ = getToken(lex).(*commentToken) }).NotTo(Panic())
		})

		It("should return rest as other", func() {
			lex := newLexer(strings.NewReader("aB2_."))
			expected := []string{"a", "B", "2", "_", "."}
			for _, e := range expected {
				token := getToken(lex)
				Expect(token.(*otherToken).value).To(Equal(e))
			}
		})
	})
})

package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"strings"
)

func checkEat(p *parser, ty tokenType) {
	t, e := p.eat(ty)
	Expect(e).To(BeNil())
	Expect(t.getType()).To(Equal(ty))
}

var _ = Describe("Parser", func() {

	It("should advance", func() {
		pars := newParser(strings.NewReader("a"))
		Expect(pars.advance()).NotTo(BeNil())
		Expect(pars.advance()).To(BeNil())
	})

	It("should eat tokens", func() {
		pars := newParser(strings.NewReader("[a]\na=b"))
		pars.advance()
		expected := []tokenType{symbolTokType, otherTokType,
			symbolTokType, newLineTokType, otherTokType, sepTokType, otherTokType}
		for _, exp := range expected {
			checkEat(pars, exp)
		}
	})

	Describe("parseSection", func() {
		It("should work with valid sections", func() {
			pars := newParser(strings.NewReader("[foo]"))
			pars.advance()
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("foo"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should transform to lower case by default", func() {
			pars := newParser(strings.NewReader("[FOO]"))
			pars.advance()
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("foo"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should not transform when option given", func() {
			pars := newParserWithOptions(strings.NewReader("[FOO]"), "[A-Za-z][A-Za-z0-9_]+", false, []byte{'='}, []byte{';'})
			pars.advance()
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("FOO"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should fail on bad section name", func() {
			pars := newParser(strings.NewReader("[f$$$o]"))
			pars.advance()
			_, err := pars.parseSection()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("Bad section name"))
		})
	})
})

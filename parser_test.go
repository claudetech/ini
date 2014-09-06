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

	Describe("advance", func() {
		It("should advance", func() {
			pars := newParser(strings.NewReader("ab"))
			Expect(pars.currentToken).NotTo(BeNil())
			Expect(pars.advance()).NotTo(BeNil())
			Expect(pars.advance()).To(BeNil())
		})
	})

	Describe("eat", func() {
		It("should eat tokens", func() {
			pars := newParser(strings.NewReader("[a]\na=b"))
			expected := []tokenType{symbolTokType, otherTokType,
				symbolTokType, newLineTokType, otherTokType, sepTokType, otherTokType}
			for _, exp := range expected {
				checkEat(pars, exp)
			}
		})
	})

	Describe("parseSection", func() {
		It("should work with valid sections", func() {
			pars := newParser(strings.NewReader("[foo]"))
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("foo"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should transform to lower case by default", func() {
			pars := newParser(strings.NewReader("[FOO]"))
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("foo"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should not transform when option given", func() {
			pars := newParserWithOptions(strings.NewReader("[FOO]"), "[A-Za-z][A-Za-z0-9_]+", false, []byte{'='}, []byte{';'})
			section, err := pars.parseSection()
			Expect(err).To(BeNil())
			Expect(section).To(Equal("FOO"))
			Expect(pars.advance()).To(BeNil())
		})

		It("should fail on bad key name", func() {
			pars := newParser(strings.NewReader("[f$$$o]"))
			_, err := pars.parseSection()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("Bad key name"))
		})
	})

	Describe("parseValue", func() {
		It("should parse normal values", func() {
			pars := newParser(strings.NewReader("my value\n"))
			v, err := pars.parseValue()
			Expect(err).To(BeNil())
			Expect(v).To(Equal("my value"))
			_, err = pars.eat(newLineTokType)
			Expect(err).To(BeNil())
		})
	})

	Describe("skipSpaces", func() {
		It("should skip all spaces", func() {
			pars := newParser(strings.NewReader("    a"))
			pars.skipSpaces()
			Expect(pars.currentToken.getType()).To(Equal(otherTokType))
		})
	})

	Describe("parseAssignment", func() {
		It("should parse valid assignments", func() {
			pars := newParser(strings.NewReader("foo = bar"))
			ident, value, err := pars.parseAssignment()
			Expect(err).To(BeNil())
			Expect(ident).To(Equal("foo"))
			Expect(value).To(Equal("bar"))
		})
	})

	Describe("parseLine", func() {
		It("should ignore blank lines", func() {
			pars := newParser(strings.NewReader("   \na"))
			err := pars.parseLine()
			Expect(err).To(BeNil())
			Expect(pars.currentToken.getType()).To(Equal(otherTokType))
		})

		It("should ignore comments lines", func() {
			pars := newParser(strings.NewReader(" ;mycomment  \na"))
			err := pars.parseLine()
			Expect(err).To(BeNil())
			Expect(pars.currentToken.getType()).To(Equal(otherTokType))
		})

		It("should handle sections", func() {
			pars := newParser(strings.NewReader("  [my_section]\na"))
			err := pars.parseLine()
			Expect(err).To(BeNil())
			_, ok := pars.currentConfig["my_section"]
			Expect(ok).To(BeTrue())
			Expect(pars.currentToken.getType()).To(Equal(otherTokType))
		})

		It("should handle assignments", func() {
			pars := newParser(strings.NewReader("  [my_section]\n foo = bar"))
			err := pars.parseLine()
			Expect(err).To(BeNil())
			_, ok := pars.currentConfig["my_section"]
			Expect(ok).To(BeTrue())
			err = pars.parseLine()
			Expect(err).To(BeNil())
			_, ok = pars.currentConfig["my_section"]["foo"]
			Expect(ok).To(BeTrue())
			Expect(pars.currentConfig["my_section"]["foo"]).To(Equal("bar"))
		})
	})

	Describe("parseConfig", func() {
		It("should parse normal config files", func() {
			conf := `
			; my config files

			[my_section]
			foo = bar ; very important

			[other_section]

			baz=qux
			`
			pars := newParser(strings.NewReader(conf))
			err := pars.parseConfig()
			Expect(err).To(BeNil())
			_, ok := pars.currentConfig["my_section"]
			Expect(ok).To(BeTrue())
			_, ok = pars.currentConfig["other_section"]
			Expect(ok).To(BeTrue())
			_, ok = pars.currentConfig["my_section"]["foo"]
			Expect(ok).To(BeTrue())
			_, ok = pars.currentConfig["other_section"]["baz"]
			Expect(ok).To(BeTrue())
			Expect(pars.currentConfig["my_section"]["foo"]).To(Equal("bar"))
			Expect(pars.currentConfig["other_section"]["baz"]).To(Equal("qux"))
		})
	})
})

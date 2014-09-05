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
		expected := []tokenType{symbolTokType, otherTokType,
			symbolTokType, newLineTokType, otherTokType, sepTokType, otherTokType}
		for _, exp := range expected {
			checkEat(pars, exp)
		}
	})

})

package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"strings"
)

var _ = Describe("Decoder", func() {
	Describe("Decode", func() {
		It("should decode to map", func() {
			config := `
			[section]
			foo=bar
			`
			d := NewDecoder(strings.NewReader(config))
			var conf map[string]map[string]string
			err := d.Decode(&conf)
			Expect(err).To(BeNil())
			_, ok := conf["section"]
			Expect(ok).To(BeTrue())
			v, ok := conf["section"]["foo"]
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal("bar"))
		})

		It("should parse complex configs", func() {
			file, err := os.Open("./test_data/php.ini")
			Expect(err).To(BeNil())
			d := NewDecoder(file)
			d.IdRegexp("[a-z][a-z0-9_ ]+")
			var conf map[string]map[string]string
			err = d.Decode(&conf)
			Expect(err).To(BeNil())
			php, ok := conf["php"]
			Expect(ok).To(BeTrue())
			engine, ok := php["engine"]
			Expect(ok).To(BeTrue())
			Expect(engine).To(Equal("On"))
		})
	})
})

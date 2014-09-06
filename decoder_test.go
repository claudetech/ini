package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"strings"
)

type conf struct {
	Section sectionConf
}

type sectionConf struct {
	Foo string
}

var _ = Describe("Decoder", func() {
	Describe("Decode", func() {
		config := `
			[section]
			foo=bar
		`
		It("should decode to map", func() {
			d := NewDecoder(strings.NewReader(config))
			var c map[string]map[string]string
			err := d.Decode(&c)
			Expect(err).To(BeNil())
			_, ok := c["section"]
			Expect(ok).To(BeTrue())
			v, ok := c["section"]["foo"]
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal("bar"))
		})

		It("should decode to struct", func() {
			d := NewDecoder(strings.NewReader(config))
			var c conf
			err := d.Decode(&c)
			Expect(err).To(BeNil())
			Expect(c.Section).NotTo(BeNil())
			Expect(c.Section.Foo).To(Equal("bar"))
		})

		It("should parse simple files", func() {
			var c map[string]map[string]string
			err := DecodeFile("./test_data/simple.ini", &c)
			Expect(err).To(BeNil())
			_, ok := c["section"]
			Expect(ok).To(BeTrue())
			v, ok := c["section"]["foo"]
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal("bar"))
		})

		It("should parse complex configs", func() {
			file, err := os.Open("./test_data/php.ini")
			defer file.Close()
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

package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
)

var _ = Describe("Encoder", func() {
	It("should encode sections", func() {
		c := new(bytes.Buffer)
		encoder := NewEncoder(c)
		err := encoder.writeSection("section", map[string]string{"foo": "bar"})
		Expect(err).To(BeNil())
		Expect(c.String()).To(Equal("[section]\nfoo = bar\n"))
	})

	It("should encode config", func() {
		c := new(bytes.Buffer)
		encoder := NewEncoder(c)
		err := encoder.Encode(map[string]map[string]string{"section": map[string]string{"foo": "bar"}})
		Expect(err).To(BeNil())
		Expect(c.String()).To(Equal("[section]\nfoo = bar\n"))
	})
})

package ini

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIni(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ini Suite")
}

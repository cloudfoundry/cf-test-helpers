package commandreporter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommandReporter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Reporter Suite")
}

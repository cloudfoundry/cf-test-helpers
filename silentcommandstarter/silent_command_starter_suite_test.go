package silentcommandstarter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSilentCommandStarter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SilentCommandStarter Suite")
}

package commandstarter_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommandStarter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Starter Suite")
}

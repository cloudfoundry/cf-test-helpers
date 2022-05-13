package commandreporter_color_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommandReporter(t *testing.T) {
	_, reportConfiguration := GinkgoConfiguration()
	reportConfiguration.NoColor = false
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Reporter Color Suite", reportConfiguration)
}

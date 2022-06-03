package commandreporter_no_color_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCommandReporter(t *testing.T) {
	_, reportConfiguration := GinkgoConfiguration()
	reportConfiguration.NoColor = true
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Reporter No Color Suite", reportConfiguration)
}

package formatter_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/formatter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CliErrorMessage", func() {
	It("returns a nice error message that includes the command that was run", func() {
		args := []string{"cf", "push"}

		argsAsString := strings.Join(args, " ")
		helpfulMessage := fmt.Sprintf("\n>>> [ %s ] exited with an error \n", argsAsString)
		Expect(formatter.CliErrorMessage(args)).To(Equal(helpfulMessage))
	})

	It("redacts args for the `cf auth` command", func() {
		args := []string{"cf", "auth", "super", "secret"}

		argsAsString := strings.Join(args[:2], " ")
		helpfulMessage := fmt.Sprintf("\n>>> [ %s ] exited with an error \n", argsAsString)
		Expect(formatter.CliErrorMessage(args)).To(Equal(helpfulMessage))
	})
})

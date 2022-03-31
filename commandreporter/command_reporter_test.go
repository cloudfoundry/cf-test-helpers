package commandreporter_test

import (
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/commandreporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("CommandReporter", func() {

	Describe("#Report", func() {
		var reporter *commandreporter.CommandReporter
		var writer *gbytes.Buffer
		var t time.Time
		var timestampRegex string
		BeforeEach(func() {
			writer = gbytes.NewBuffer()
			reporter = commandreporter.NewCommandReporter()
			t = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
			timestampRegex = "\\[2009-11-10 23:00:00.00 \\(UTC\\)\\]>"
		})

		It("prints the timestamp and command in green", func() {
			cmd := exec.Command("executable", "arg1", "arg2")
			reporter.Report(t, cmd)

			lineStart := "^\n"
			greenStart := "\\x1b\\[32m"
			greenEnd := "\\x1b\\[0m"
			lineEnd := "\n$"
			Expect(writer).To(gbytes.Say("%s%s%s executable arg1 arg2 %s%s", lineStart, greenStart, timestampRegex, greenEnd, lineEnd))
		})

		Context("when NoColor is specified", func() {

			It("does not print color", func() {
				cmd := exec.Command("executable", "arg1", "arg2")
				reporter.Report(t, cmd)

				lineStart := "^\n"
				lineEnd := "\n$"
				Expect(writer).To(gbytes.Say("%s%s executable arg1 arg2 %s", lineStart, timestampRegex, lineEnd))

			})
		})
	})
})

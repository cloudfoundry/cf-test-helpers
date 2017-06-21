package internal_test

import (
	"time"

	"os/exec"

	"bytes"

	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedactingReporter", func() {
	Describe("Report", func() {
		var (
			startTime time.Time
			cmd       *exec.Cmd
			buffer    *bytes.Buffer

			reporter internal.Reporter
		)

		BeforeEach(func() {
			buffer = &bytes.Buffer{}
		})

		It("prints the time", func() {
			cmd = exec.Command("some-command", "with", "args")
			reporter = internal.NewRedactingReporter(buffer)

			reporter.Report(startTime, cmd)

			Expect(buffer.String()).To(ContainSubstring("[0001-01-01 00:00:00.00 (UTC)]>"))
		})

		Context("Does not redact", func() {
			It("if no redactees specified", func() {
				cmd = exec.Command("whatever", "create-org", "foo")
				reporter = internal.NewRedactingReporter(buffer)

				reporter.Report(startTime, cmd)

				Expect(buffer.String()).To(ContainSubstring("whatever create-org foo"))
			})

			It("if no matching redactees", func() {
				cmd = exec.Command("blah", "important")
				reporter = internal.NewRedactingReporter(buffer, "feh", "meh")

				reporter.Report(startTime, cmd)

				Expect(buffer.String()).NotTo(ContainSubstring("[REDACTED]"))
				Expect(buffer.String()).To(ContainSubstring("blah important"))
			})
		})

		Context("Redacts", func() {
			It("one value", func() {
				cmd = exec.Command("blah", "important")
				reporter = internal.NewRedactingReporter(buffer, "important")

				reporter.Report(startTime, cmd)

				Expect(buffer.String()).NotTo(ContainSubstring("important"))
				Expect(buffer.String()).To(ContainSubstring("blah [REDACTED]"))
			})

			It("multiple values", func() {
				cmd = exec.Command("command", "sensitive", "secret", "other")
				reporter = internal.NewRedactingReporter(buffer, "sensitive", "secret")

				reporter.Report(startTime, cmd)

				Expect(buffer.String()).NotTo(ContainSubstring("sensitive"))
				Expect(buffer.String()).NotTo(ContainSubstring("secret"))
				Expect(buffer.String()).To(ContainSubstring("command [REDACTED] [REDACTED] other"))
			})
		})
	})
})

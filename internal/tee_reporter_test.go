package internal_test

import (
	"os/exec"
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type reportOnly struct{ reports int }

func (r *reportOnly) Report(time.Time, *exec.Cmd) { r.reports++ }

type reportAndComplete struct {
	reports   int
	completes int
}

func (r *reportAndComplete) Report(time.Time, *exec.Cmd) { r.reports++ }
func (r *reportAndComplete) ReportCompletion(*exec.Cmd, time.Duration, int) {
	r.completes++
}

var _ = Describe("TeeReporter", func() {
	It("fans Report to every reporter but ReportCompletion only to completion-capable ones", func() {
		a := &reportOnly{} // no ReportCompletion — must be skipped, no panic
		b := &reportAndComplete{}
		tee := internal.NewTeeReporter(a, b)

		tee.Report(time.Now(), exec.Command("cf", "app"))
		Expect(a.reports).To(Equal(1))
		Expect(b.reports).To(Equal(1))

		tee.ReportCompletion(exec.Command("cf", "app"), 10*time.Millisecond, 0)
		Expect(b.completes).To(Equal(1))
	})
})

package internal_test

import (
	"io"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cf-test-helpers/v2/internal"
)

type observedBaseReporter struct{}

func (observedBaseReporter) Report(time.Time, *exec.Cmd) {}

var _ = Describe("ObservedReporter", func() {
	AfterEach(func() {
		internal.UnregisterObserver()
	})

	It("returns the base reporter unchanged when no observer is registered", func() {
		base := observedBaseReporter{}
		Expect(internal.ObservedReporter(base)).To(BeIdenticalTo(internal.Reporter(base)))
	})

	Context("when an observer is registered", func() {
		var observer *recordingObserver

		BeforeEach(func() {
			observer = &recordingObserver{}
			internal.RegisterObserver(observer)
		})

		It("redacts the observed identity when the base carries a redactor", func() {
			base := internal.NewRedactingReporter(io.Discard, internal.NewRedactor("secret"))
			got := internal.ObservedReporter(base)

			cmd := exec.Command("cf", "auth", "user", "secret")
			got.Report(time.Now(), cmd)
			got.(internal.CompletionReporter).ReportCompletion(cmd, time.Second, 0)

			Expect(observer.startedArgs).To(Equal("cf auth user [REDACTED]"))
			Expect(observer.completedArgs).To(Equal("cf auth user [REDACTED]"))
		})
	})
})

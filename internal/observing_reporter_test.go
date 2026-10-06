package internal_test

import (
	"os/exec"
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type recordingObserver struct {
	startedArgs   string
	startedTime   time.Time
	completedArgs string
	elapsed       time.Duration
	exitCode      int
	startedCalls  int
	doneCalls     int
}

func (o *recordingObserver) CommandStarted(redactedArgs string, startTime time.Time) {
	o.startedArgs = redactedArgs
	o.startedTime = startTime
	o.startedCalls++
}

func (o *recordingObserver) CommandCompleted(redactedArgs string, elapsed time.Duration, exitCode int) {
	o.completedArgs = redactedArgs
	o.elapsed = elapsed
	o.exitCode = exitCode
	o.doneCalls++
}

var _ = Describe("ObservingReporter", func() {
	var (
		observer *recordingObserver
		reporter *internal.ObservingReporter
	)

	BeforeEach(func() {
		observer = &recordingObserver{}
		// redactor redacts the multi-word secret "top secret"
		reporter = internal.NewObservingReporter(observer, internal.NewRedactor("top secret"))
	})

	It("forwards a redacted start observation", func() {
		start := time.Now()
		cmd := exec.Command("cf", "auth", "top", "secret")

		reporter.Report(start, cmd)

		Expect(observer.startedCalls).To(Equal(1))
		Expect(observer.startedArgs).To(Equal("cf auth [REDACTED]"))
		Expect(observer.startedArgs).NotTo(ContainSubstring("top secret"))
		Expect(observer.startedTime).To(Equal(start))
	})

	It("forwards a redacted completion observation with elapsed and exit code", func() {
		cmd := exec.Command("cf", "auth", "top", "secret")

		reporter.ReportCompletion(cmd, 250*time.Millisecond, 7)

		Expect(observer.doneCalls).To(Equal(1))
		Expect(observer.completedArgs).To(Equal("cf auth [REDACTED]"))
		Expect(observer.completedArgs).NotTo(ContainSubstring("top secret"))
		Expect(observer.elapsed).To(Equal(250 * time.Millisecond))
		Expect(observer.exitCode).To(Equal(7))
	})
})

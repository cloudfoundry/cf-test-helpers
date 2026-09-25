package internal_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-test-helpers/v2/commandreporter"
	"github.com/cloudfoundry/cf-test-helpers/v2/internal"
	"github.com/cloudfoundry/cf-test-helpers/v2/internal/fakes"
)

// syncObserver is a concurrency-safe CommandObserver for asserting on the
// start observation (synchronous) and the completion observation (delivered on
// the starter's goroutine) under Eventually.
type syncObserver struct {
	mu          sync.Mutex
	starts      []string
	completions []string
}

func (o *syncObserver) CommandStarted(redactedArgs string, _ time.Time) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.starts = append(o.starts, redactedArgs)
}

func (o *syncObserver) CommandCompleted(redactedArgs string, _ time.Duration, _ int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.completions = append(o.completions, redactedArgs)
}

func (o *syncObserver) Starts() []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]string(nil), o.starts...)
}

func (o *syncObserver) Completions() []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]string(nil), o.completions...)
}

func (o *syncObserver) Snapshot() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return strings.Join(append(o.starts, o.completions...), " ")
}

var _ = Describe("cf", func() {
	var (
		starter *fakes.FakeCmdStarter
	)
	BeforeEach(func() {
		starter = fakes.NewFakeCmdStarter()
	})
	AfterEach(func() {
		internal.UnregisterObserver()
	})

	Describe("Cf", func() {
		It("uses a default reporter", func() {
			Eventually(internal.Cf(starter, "app", "my-app"), 1*time.Second).Should(Exit(0))
			Expect(starter.CalledWith[0].Reporter).To(BeAssignableToTypeOf(commandreporter.NewCommandReporter()))
		})
	})

	Describe("CfWithCustomReporter", func() {
		var (
			buffer            *bytes.Buffer
			fakeRedactor      *fakes.FakeRedactor
			redactingReporter internal.Reporter
		)
		BeforeEach(func() {
			buffer = &bytes.Buffer{}
			fakeRedactor = &fakes.FakeRedactor{}
			redactingReporter = internal.NewRedactingReporter(buffer, fakeRedactor)
		})

		It("calls the cf cli with the correct command and args", func() {
			Eventually(internal.CfWithCustomReporter(starter, redactingReporter, "app", "my-app"), 1*time.Second).Should(Exit(0))

			Expect(starter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(starter.CalledWith[0].Args).To(Equal([]string{"app", "my-app"}))
		})

		It("uses the specified reporter", func() {
			Eventually(internal.CfWithCustomReporter(starter, redactingReporter, "app", "my-app"), 1*time.Second).Should(Exit(0))
			Expect(starter.CalledWith[0].Reporter).To(BeAssignableToTypeOf(internal.NewRedactingReporter(io.Discard, fakeRedactor)))
		})

		Context("when the exit code is non-zero", func() {
			BeforeEach(func() {
				starter.ToReturn[0].ExitCode = 42
			})

			It("returns the exit code anyway", func() {
				Eventually(internal.CfWithCustomReporter(starter, redactingReporter, "app", "my-app"), 1*time.Second).Should(Exit(42))
			})
		})

		Context("when there is an error", func() {
			BeforeEach(func() {
				starter.ToReturn[0].Err = fmt.Errorf("failing now")
			})

			It("panics", func() {
				Expect(func() {
					internal.CfWithCustomReporter(starter, redactingReporter, "fail")
				}).To(Panic())
			})
		})
	})

	Describe("observer seam", func() {
		It("leaves the base reporter untouched when no observer is registered", func() {
			Eventually(internal.Cf(starter, "app", "my-app"), 1*time.Second).Should(Exit(0))
			Expect(starter.CalledWith[0].Reporter).To(BeAssignableToTypeOf(commandreporter.NewCommandReporter()))
			Expect(starter.CalledWith[0].Reporter).NotTo(BeAssignableToTypeOf(&internal.TeeReporter{}))
		})

		It("delivers redacted start and completion to a registered observer", func() {
			observer := &syncObserver{}
			internal.RegisterObserver(observer)
			redactingReporter := internal.NewRedactingReporter(io.Discard, internal.NewRedactor("secret"))

			session := internal.CfWithCustomReporter(starter, redactingReporter, "auth", "user", "secret")
			Eventually(session, 1*time.Second).Should(Exit(0))

			Eventually(observer.Starts, 1*time.Second).Should(ConsistOf("cf auth user [REDACTED]"))
			Eventually(observer.Completions, 1*time.Second).Should(ConsistOf("cf auth user [REDACTED]"))
			Expect(observer.Snapshot()).NotTo(ContainSubstring("secret"))
		})
	})
})

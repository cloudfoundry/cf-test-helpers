package cf_test

import (
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry/cf-test-helpers/v2/cf"
	"github.com/cloudfoundry/cf-test-helpers/v2/internal"
	"github.com/cloudfoundry/cf-test-helpers/v2/internal/fakes"
)

type fakeObserver struct {
	mu          sync.Mutex
	starts      []string
	completions []string
}

func (o *fakeObserver) CommandStarted(redactedArgs string, _ time.Time) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.starts = append(o.starts, redactedArgs)
}

func (o *fakeObserver) CommandCompleted(redactedArgs string, _ time.Duration, _ int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.completions = append(o.completions, redactedArgs)
}

func (o *fakeObserver) Starts() []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]string(nil), o.starts...)
}

func (o *fakeObserver) Completions() []string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return append([]string(nil), o.completions...)
}

// These specs cover the two things the internal seam tests cannot: that the
// flow works through the exported cf.RegisterObserver surface, and that
// registration is race-clean under concurrent callers. The seam's redaction
// and opt-out behavior are proven in internal/cf_test.go.
var _ = Describe("cf command observer (end to end)", func() {
	var observer *fakeObserver

	BeforeEach(func() {
		observer = &fakeObserver{}
	})

	AfterEach(func() {
		cf.UnregisterObserver()
	})

	It("emits start and completion through the exported cf.RegisterObserver surface", func() {
		cf.RegisterObserver(observer)
		starter := fakes.NewFakeCmdStarter()

		session := internal.Cf(starter, "app", "my-app")
		Eventually(session, 1*time.Second).Should(Exit(0))

		Eventually(observer.Starts, 1*time.Second).Should(ConsistOf("cf app my-app"))
		Eventually(observer.Completions, 1*time.Second).Should(ConsistOf("cf app my-app"))
	})

	It("stays race-clean under concurrent callers and registration churn", func() {
		cf.RegisterObserver(observer)

		const workers = 8
		var wg sync.WaitGroup
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				// Each goroutine gets its own starter: FakeCmdStarter is not
				// concurrency-safe, so sharing one would race the fake rather
				// than the code under test.
				starter := fakes.NewFakeCmdStarter()
				session := internal.Cf(starter, "app", "my-app")
				Eventually(session, 1*time.Second).Should(Exit(0))
			}()
		}

		for i := 0; i < workers; i++ {
			cf.RegisterObserver(observer)
			cf.UnregisterObserver()
		}

		wg.Wait()
	})
})

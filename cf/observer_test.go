package cf

// NOTE: white-box test in package cf so it can read the unexported
// registeredObserver() directly, rather than exporting a test-only accessor
// from production code.

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type noopObserver struct{}

func (noopObserver) CommandStarted(string, time.Time)           {}
func (noopObserver) CommandCompleted(string, time.Duration, int) {}

var _ = Describe("Observer registration", func() {
	AfterEach(func() {
		UnregisterObserver()
	})

	It("registers and unregisters a cf.Observer through the exported surface", func() {
		var obs Observer = noopObserver{}
		Expect(registeredObserver()).To(BeNil())

		RegisterObserver(obs)
		Expect(registeredObserver()).NotTo(BeNil())

		UnregisterObserver()
		Expect(registeredObserver()).To(BeNil())
	})
})

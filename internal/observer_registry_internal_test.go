package internal

import (
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type registryNoopObserver struct{}

func (registryNoopObserver) CommandStarted(string, time.Time)            {}
func (registryNoopObserver) CommandCompleted(string, time.Duration, int) {}

var _ = ginkgo.Describe("observer registry", func() {
	ginkgo.AfterEach(func() {
		UnregisterObserver()
	})

	ginkgo.It("registers and unregisters an observer", func() {
		obs := registryNoopObserver{}
		RegisterObserver(obs)
		gomega.Expect(RegisteredObserver()).NotTo(gomega.BeNil())

		UnregisterObserver()
		gomega.Expect(RegisteredObserver()).To(gomega.BeNil())
	})
})

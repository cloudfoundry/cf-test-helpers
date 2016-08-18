package cfinternal_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf/internal"
)

var _ = Describe("Cf", func() {
	It("calls the cf cli with the correct command and args", func() {
		starter := new(fakeStarter)
		starter.toReturn.exitCode = 42

		Eventually(cfinternal.Cf(starter, "app", "my-app"), 1*time.Second).Should(Exit(42))

		Expect(starter.calledWith.executable).To(Equal("cf"))
		Expect(starter.calledWith.args).To(Equal([]string{"app", "my-app"}))
	})

	Context("when there is an error", func() {
		It("panics", func() {
			starter := new(fakeStarter)
			starter.toReturn.err = fmt.Errorf("failing now")
			Expect(func() {
				cfinternal.Cf(starter, "fail")
			}).To(Panic())
		})
	})
})

package cfinternal_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf/internal"
	"github.com/cloudfoundry-incubator/cf-test-helpers/commandstarter"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal/fakes"
)

var _ = Describe("Cf", func() {
	It("calls the cf cli with the correct command and args", func() {
		starter := new(fakes.FakeCmdStarter)
		starter.ToReturn.ExitCode = 42

		Eventually(cfinternal.Cf(starter, "app", "my-app"), 1*time.Second).Should(Exit(42))

		Expect(starter.CalledWith.Executable).To(Equal("cf"))
		Expect(starter.CalledWith.Args).To(Equal([]string{"app", "my-app"}))
	})

	It("uses a default reporter", func() {
		starter := new(fakes.FakeCmdStarter)
		Eventually(cfinternal.Cf(starter, "app", "my-app"), 1*time.Second).Should(Exit(0))
		Expect(starter.CalledWith.Reporter).To(BeAssignableToTypeOf(commandstarter.NewDefaultReporter()))
	})

	Context("when there is an error", func() {
		It("panics", func() {
			starter := new(fakes.FakeCmdStarter)
			starter.ToReturn.Err = fmt.Errorf("failing now")
			Expect(func() {
				cfinternal.Cf(starter, "fail")
			}).To(Panic())
		})
	})
})

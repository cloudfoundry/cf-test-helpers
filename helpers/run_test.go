package helpers_test

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Run", func() {
	var cmdTimeout time.Duration
	BeforeEach(func() {
		cmdTimeout = 30 * time.Second
	})

	It("runs the given command in a cmdtest Session", func() {
		session := helpers.Run("bash", "-c", "echo hi out; echo hi err 1>&2; exit 42")

		session.Wait(cmdTimeout)
		Expect(session).To(Exit(42))
		Expect(session.Out).To(Say("hi out"))
		Expect(session.Err).To(Say("hi err"))
	})
})

var _ = Describe("cmdhelpers", func() {
	Describe("Run with defaults", func() {
		var cmdTimeout time.Duration
		var sleepTimeInSeconds int // this is an int instead of a time.Duration so we can render it easily in bash strings
		Context("when the command succeeds before the timeout", func() {
			BeforeEach(func() {
				cmdTimeout = 30 * time.Second
				sleepTimeInSeconds = 1
			})

			It("does nothing", func() {
				failures := InterceptGomegaFailures(func() {
					session := helpers.Run("bash", "-c", fmt.Sprintf("sleep %d; echo hi out; echo hi err 1>&2; exit 0", sleepTimeInSeconds))
					session.Wait(cmdTimeout)
				})
				Expect(failures).To(BeEmpty())
			})
		})

		Context("when the command takes longer than timeout", func() {
			BeforeEach(func() {
				cmdTimeout = 1 * time.Second
				sleepTimeInSeconds = 10
			})

			It("fails the ginkgo test", func() {
				failures := InterceptGomegaFailures(func() {
					session := helpers.Run("bash", "-c", fmt.Sprintf("sleep %d; echo hi out; echo hi err 1>&2; exit 0", sleepTimeInSeconds))
					session.Wait(cmdTimeout)
				})
				Expect(failures).NotTo(BeEmpty())
			})
		})
	})
})

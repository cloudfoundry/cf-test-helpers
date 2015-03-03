package runner_test

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

const CMD_TIMEOUT = 30 * time.Second

var _ = Describe("Run", func() {
	It("runs the given command in a cmdtest Session", func() {
		session := runner.Run("bash", "-c", "echo hi out; echo hi err 1>&2; exit 42").Wait(CMD_TIMEOUT)
		Expect(session).To(Exit(42))
		Expect(session.Out).To(Say("hi out"))
		Expect(session.Err).To(Say("hi err"))
	})
})

var _ = Describe("Curl", func() {
	It("outputs the body of the given URL", func() {
		session := runner.Curl("-I", "http://example.com").Wait(CMD_TIMEOUT)
		Expect(session).To(Exit(0))
		Expect(session.Out).To(Say("HTTP/1.1 200 OK"))
	})
})

var _ = Describe("ExecWithTimeout", func() {
	It("does nothing when the command succeeds before the timeout", func() {
		failures := InterceptGomegaFailures(func() {
			session := runner.Run("bash", "-c", "echo hi out; echo hi err 1>&2; exit 0")
			runner.ExecWithTimeout(session, CMD_TIMEOUT)
		})
		Expect(failures).To(BeEmpty())
	})

	It("expects the command to not fail", func() {
		failures := InterceptGomegaFailures(func() {
			session := runner.Run("bash", "-c", "echo hi out; echo hi err 1>&2; exit 42")
			runner.ExecWithTimeout(session, CMD_TIMEOUT)
		})
		Expect(failures[0]).To(MatchRegexp(
			"Failed executing command \\(exit 42\\):\nCommand: %s\n\n\\[stdout\\]:\n%s\n\n\\[stderr\\]:\n%s",
			"bash -c echo hi out; echo hi err 1>&2; exit 42",
			"hi out\n",
			"hi err\n",
		))
	})

	It("expects the command to not time out", func() {
		failures := InterceptGomegaFailures(func() {
			session := runner.Run("bash", "-c", "echo hi out; echo hi err 1>&2; sleep 1")
			runner.ExecWithTimeout(session, 100*time.Millisecond)
		})
		Expect(failures[0]).To(MatchRegexp(
			"Timed out executing command \\(100ms\\):\nCommand: %s\n\n\\[stdout\\]:\n%s\n\n\\[stderr\\]:\n%s",
			"bash -c echo hi out; echo hi err 1>&2; sleep 1",
			"hi out\n",
			"hi err\n",
		))
	})
})

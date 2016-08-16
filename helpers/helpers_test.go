package helpers_test

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Curl", func() {
	var cmdTimeout time.Duration
	BeforeEach(func() {
		cmdTimeout = 30 * time.Second
	})

	It("outputs the body of the given URL", func() {
		session := helpers.Curl("-I", "http://example.com")

		session.Wait(cmdTimeout)
		Expect(session).To(Exit(0))
		Expect(session.Out).To(Say("HTTP/1.1 200 OK"))
	})
})

var _ = Describe("Helpers", func() {
	It("builds", func() {})
})

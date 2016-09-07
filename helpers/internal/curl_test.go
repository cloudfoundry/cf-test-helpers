package helpersinternal_test

import (
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers/internal"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Curl", func() {
	var cmdTimeout time.Duration
	BeforeEach(func() {
		cmdTimeout = 30 * time.Second
	})

	It("outputs the body of the given URL", func() {
		starter := new(fakes.FakeCmdStarter)
		starter.ToReturn.Output = "HTTP/1.1 200 OK"

		session := helpersinternal.Curl(starter, false, "-I", "http://example.com")

		session.Wait(cmdTimeout)
		Expect(session).To(gexec.Exit(0))
		Expect(session.Out).To(Say("HTTP/1.1 200 OK"))
		Expect(starter.CalledWith[0].Executable).To(Equal("curl"))
		Expect(starter.CalledWith[0].Args).To(ConsistOf("-I", "-s", "http://example.com"))
	})

	It("panics when the starter returns an error", func() {
		starter := new(fakes.FakeCmdStarter)
		starter.ToReturn.Err = fmt.Errorf("error")

		Expect(func() {
			helpersinternal.Curl(starter, false, "-I", "http://example.com")
		}).To(Panic())
	})
})

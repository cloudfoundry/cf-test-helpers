package cfinternal_test

import (
	"fmt"
	"os/exec"
	"time"

	. "github.com/cloudfoundry-incubator/cf-test-helpers/cf/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

type fakeStarter struct {
	calledWith struct {
		executable string
		args       []string
	}
	toReturn struct {
		output    string
		err       error
		exitCode  int
		sleepTime int
	}
}

func (s *fakeStarter) Start(executable string, args ...string) (*gexec.Session, error) {
	s.calledWith.executable = executable
	s.calledWith.args = args

	// Default return values
	if s.toReturn.output == "" {
		s.toReturn.output = `\{\}`
	}
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"echo %s; sleep %d; exit %d",
			s.toReturn.output,
			s.toReturn.sleepTime,
			s.toReturn.exitCode,
		),
	)
	session, _ := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	return session, s.toReturn.err
}

type genericResource struct {
	Foo string
}

var _ = Describe("ApiRequest", func() {
	var starter *fakeStarter
	var timeout time.Duration
	BeforeEach(func() {
		starter = new(fakeStarter)
		starter.toReturn.output = `\{\"foo\": \"bar\"\}`
		timeout = 1 * time.Second
	})

	It("sends the request to the current CF target", func() {
		var response genericResource
		ApiRequest(starter, "GET", "/v2/info", &response, timeout, "some", "data")
		Expect(starter.calledWith.executable).To(Equal("cf"))
		Expect(starter.calledWith.args).To(Equal([]string{"curl", "/v2/info", "-X", "GET", "-d", "somedata"}))

		Expect(response.Foo).To(Equal("bar"))
	})

	Context("when command starter returns a nil response", func() {
		It("does not fail", func() {
			failures := InterceptGomegaFailures(func() {
				ApiRequest(starter, "GET", "/v2/info", nil, timeout)
			})
			Expect(failures).To(BeEmpty())
		})
	})

	Context("when request data is empty", func() {
		It("doesn't include a -d argument", func() {
			ApiRequest(starter, "GET", "/v2/info", nil, timeout)
			Expect(starter.calledWith.args).NotTo(ContainElement("-d"))
		})
	})

	Context("when there is an error from the starter", func() {
		BeforeEach(func() {
			starter.toReturn.err = fmt.Errorf("something went wrong")
		})

		It("fails with a ginkgo error", func() {
			failures := InterceptGomegaFailures(func() {
				ApiRequest(starter, "GET", "/v2/info", nil, timeout)
			})
			Expect(failures).To(ContainElement(MatchRegexp("something went wrong\n.*not to have occurred")))
		})
	})

	Context("when the process exits with non-zero code", func() {
		BeforeEach(func() {
			starter.toReturn.exitCode = 1
		})

		It("fails with a ginkgo error", func() {
			failures := InterceptGomegaFailures(func() {
				ApiRequest(starter, "GET", "/v2/info", nil, timeout)
			})
			Expect(failures).To(ContainElement(MatchRegexp("1\n.*to match exit code:\n.*0")))
		})
	})

	Context("when the process takes too long to exit", func() {
		BeforeEach(func() {
			timeout = 1 * time.Millisecond
		})

		It("fails with a ginkgo error", func() {
			failures := InterceptGomegaFailures(func() {
				ApiRequest(starter, "GET", "/v2/info", nil, timeout)
			})
			Expect(failures).To(ContainElement(MatchRegexp("Timed out after 0.001s.\nExpected process to exit.  It did not.")))
		})
	})

	Context("when there is an error unmarshaling the api response", func() {
		BeforeEach(func() {
			starter.toReturn.output = `{{{`
		})

		It("fails with a ginkgo error", func() {
			failures := InterceptGomegaFailures(func() {
				var response genericResource
				ApiRequest(starter, "GET", "/v2/info", &response, timeout)
			})
			Expect(failures).To(ContainElement(MatchRegexp("json.SyntaxError")))
		})
	})
})

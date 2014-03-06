package cf_test

import (
	"errors"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"github.com/vito/cmdtest"
)

var _ = Describe("AsUser", func() {
	var FakeThingsToRunAsUser = func() error { return nil }
	var FakeCfCalls = [][]string{}

	var FakeCf = func(args ...string) *cmdtest.Session {
		FakeCfCalls = append(FakeCfCalls, args)
		var session, _ = cmdtest.Start(exec.Command("echo", "nothing"))
		return session
	}
	var user = cf.NewUser("FAKE_USERNAME", "FAKE_PASSWORD", "FAKE_ORG", "FAKE_SPACE")

	BeforeEach(func() {
		FakeCfCalls = [][]string{}
		cf.Cf = FakeCf
	})

	It("calls cf login", func() {
		cf.AsUser(user, FakeThingsToRunAsUser)

		Expect(FakeCfCalls[0]).To(Equal([]string{"login", "FAKE_USERNAME", "FAKE_PASSWORD"}))
	})

	It("calls the passed function", func() {
		called := false
		cf.AsUser(user, func() error { called = true; return nil })

		Ω(called).To(BeTrue())
	})

	It("calls cf logout", func() {
		cf.AsUser(user, FakeThingsToRunAsUser)

		Expect(FakeCfCalls[1]).To(Equal([]string{"logout"}))
	})

	It("logs out even if there's an error", func() {
		cf.AsUser(user, func() error { return errors.New("_") })

		Expect(FakeCfCalls[len(FakeCfCalls)-1]).To(Equal([]string{"logout"}))
	})

	Context("when the passed function returns an error", func() {
		It("returns the same error", func() {
			myError := errors.New("fake error")

			Expect(cf.AsUser(user, func() error { return myError })).To(Equal(myError))
		})
	})

	It("sets a unique CF_HOME value", func() {
		var (
			firstHome string
			secondHome string
		)

		cf.AsUser(user, func() error {
			firstHome = os.Getenv("CF_HOME")
			return nil
		})

		cf.AsUser(user, func() error {
			secondHome = os.Getenv("CF_HOME")
			return nil
		})

		Expect(firstHome).NotTo(Equal(secondHome))
	})

	It("returns CF_HOME to its original value", func() {
		os.Setenv("CF_HOME", "some-crazy-value")
		cf.AsUser(user, func() error { return nil })

		Expect(os.Getenv("CF_HOME")).To(Equal("some-crazy-value"))
	})
})

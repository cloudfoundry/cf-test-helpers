package cf_test

import (
	"errors"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"github.com/vito/cmdtest"
)

var _ = Describe("CfAsUser", func() {
	var FakeThingsToRunAsUser = func() error { return nil }
	var FakeCfCalls = [][]string{}

	var FakeCf = func(args ...string) *cmdtest.Session {
		FakeCfCalls = append(FakeCfCalls, args)
		var session, _ = cmdtest.Start(exec.Command("echo", "nothing"))
		return session
	}

	BeforeEach(func() {
		FakeCfCalls = [][]string{}
		cf.Cf = FakeCf
	})

	It("calls cf login", func(){
		cf.CfAsUser("uname", "passwd", FakeThingsToRunAsUser)

		Expect(FakeCfCalls[0]).To(Equal([]string{"login", "uname", "passwd"}))
	})

	It("calls the passed function", func(){
		called := false
		cf.CfAsUser("uname", "passwd", func() error{ called = true; return nil })

		Ω(called).To(BeTrue())
	})

	It("calls cf login", func(){
		cf.CfAsUser("_", "_", FakeThingsToRunAsUser)

		Expect(FakeCfCalls[1]).To(Equal([]string{"logout"}))
	})

	It("logs out even if there's an error", func(){
		cf.CfAsUser("_", "_", func() error { return errors.New("_") })

		Expect(FakeCfCalls[len(FakeCfCalls) - 1]).To(Equal([]string{"logout"}))
	})

	Context("when the passed function returns an error", func(){
		It("returns the same error", func(){
			myError := errors.New("fake error")

			Expect(cf.CfAsUser("_", "_", func() error { return myError })).To(Equal(myError))
		})
	})
})

package cf_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
)

var _ = Describe("NewUserContext", func() {

	var createUser = func() cf.UserContext {
		return cf.NewUserContext("FAKE_USERNAME", "FAKE_PASSWORD", "FAKE_ORG", "FAKE_SPACE")
	}

	It("returns a UserContext struct", func() {
		Expect(createUser()).To(BeAssignableToTypeOf(cf.UserContext{}))
	})

	It("sets UserContext.name", func() {
		Expect(createUser().Username).To(Equal("FAKE_USERNAME"))
	})

	It("sets UserContext.password", func() {
		Expect(createUser().Password).To(Equal("FAKE_PASSWORD"))
	})

	It("sets UserContext.org", func() {
		Expect(createUser().Org).To(Equal("FAKE_ORG"))
	})

	It("sets UserContext.space", func() {
		Expect(createUser().Space).To(Equal("FAKE_SPACE"))
	})
})

package cf_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
)

var _ = Describe("NewUser", func() {

	var createUser = func() cf.User {
		return cf.NewUser("FAKE_USERNAME", "FAKE_PASSWORD", "FAKE_ORG", "FAKE_SPACE")
	}

	It("returns a User struct", func() {
		Expect(createUser()).To(BeAssignableToTypeOf(cf.User{}))
	})

	It("sets User.name", func() {
		Expect(createUser().Username).To(Equal("FAKE_USERNAME"))
	})

	It("sets User.password", func() {
		Expect(createUser().Password).To(Equal("FAKE_PASSWORD"))
	})

	It("sets User.org", func() {
		Expect(createUser().Org).To(Equal("FAKE_ORG"))
	})

	It("sets User.space", func() {
		Expect(createUser().Space).To(Equal("FAKE_SPACE"))
	})
})

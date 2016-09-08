package workflowhelpers

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Context", func() {
	Describe("NewContext", func() {
		It("returns a ConfiguredContext struc", func() {})

		It("sets config", func() {})

		It("sets shortTimeout and longTimeout", func() {})

		It("sets regularUserPassword and regularUserUsername", func() {})

		It("sets adminUserUsername and adminUserPassword", func() {})

		It("sets the test space", func() {})
	})

	Describe("NewPersistentAppContext", func() {
		It("returns a ConfiguredContext struc", func() {})

		It("sets the test space", func() {})
	})

	Describe("NewRunawayAppTestSetup", func() {
		It("sets the test space", func() {})
	})

	Describe("Setup", func() {
		Context("Logs in as the admin", func() {
			It("creates a new user when config.UseExistingUser is false", func() {})
		})

		Context("Logs in as a regular user", func() {

		})
	})

	Describe("TearDown", func() {

		It("Logs out the regular user", func() {})

		It("Retores CfHomeDir", func() {})

		Context("Logs in as the admin", func() {
			Context("when config.ShouldKeepUser is false", func() {
				It("deletes the user", func() {
				})
			})
		})
	})

	Describe("AdminUserContext", func() {
		It("sets ApiEndpoint", func() {})

		It("sets username", func() {})

		It("sets password", func() {})

		It("leaves org name empty", func() {})

		It("leaves space name empty", func() {})

		It("sets SkipSSLValidation", func() {})
	})

	Describe("RegularUserContext", func() {
		It("sets ApiEndpoint", func() {})

		It("sets username", func() {})

		It("sets password", func() {})

		It("sets org name", func() {})

		It("sets space name", func() {})

		It("sets SkipSSLValidation", func() {})

	})
})

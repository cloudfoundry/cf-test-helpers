package internal_test

import (
	"fmt"

	"github.com/cloudfoundry/cf-test-helpers/v2/config"
	"github.com/cloudfoundry/cf-test-helpers/v2/internal/fakes"
	. "github.com/cloudfoundry/cf-test-helpers/v2/workflowhelpers/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	var cfg *config.Config

	Describe("NewTestUser", func() {
		var existingUser, existingUserPassword, userOrigin string
		var useExistingUser bool
		var configurableTestPassword string

		BeforeEach(func() {
			useExistingUser = false
			configurableTestPassword = ""
			userOrigin = ""
		})

		JustBeforeEach(func() {
			cfg = &config.Config{
				NamePrefix:               "UNIT-TESTS",
				UseExistingUser:          useExistingUser,
				ExistingUser:             existingUser,
				UserOrigin:               userOrigin,
				ExistingUserPassword:     existingUserPassword,
				ConfigurableTestPassword: configurableTestPassword,
			}
		})

		It("has a username with the name prefix and a password with default length", func() {
			user := NewTestUser(cfg, &fakes.FakeCmdStarter{})
			Expect(user.Username()).To(MatchRegexp("UNIT-TESTS-[0-9]+-USER-.*"))
			Expect(len(user.Password())).To(Equal(20))
		})

		It("has a unique random password", func() {
			password1 := NewTestUser(cfg, &fakes.FakeCmdStarter{}).Password()
			password2 := NewTestUser(cfg, &fakes.FakeCmdStarter{}).Password()
			Expect(password1).ToNot(Equal(password2))
		})

		Context("when the user origin is specified", func() {
			BeforeEach(func() {
				userOrigin = "my-test-user-origin"
			})

			It("uses the Origin", func() {
				user := NewTestUser(cfg, &fakes.FakeCmdStarter{})
				Expect(user.Origin()).To(Equal(userOrigin))
			})

		})

		Context("when the config specifies that an existing user should be used", func() {
			BeforeEach(func() {
				useExistingUser = true
				existingUser = "my-test-user"
				existingUserPassword = "my-test-password"
			})

			It("uses the ExistingUser and ExistingUserPassword", func() {
				user := NewTestUser(cfg, &fakes.FakeCmdStarter{})
				Expect(user.Username()).To(Equal(existingUser))
				Expect(user.Password()).To(Equal(existingUserPassword))
			})
		})

		Context("when the config includes a ConfigurableTestPassword", func() {
			BeforeEach(func() {
				configurableTestPassword = "pre-configured-test-password"
			})

			It("uses the a random user name and the ConfigurableTestPassword", func() {
				user := NewTestUser(cfg, &fakes.FakeCmdStarter{})
				Expect(user.Username()).To(MatchRegexp("UNIT-TESTS-[0-9]+-USER-.*"))
				Expect(user.Password()).To(Equal(configurableTestPassword))
			})
		})
	})

	Describe("NewAdminUser", func() {
		It("copies the username and password from the config", func() {
			cfg := &config.Config{AdminUser: "admin", AdminPassword: "admin-password"}
			user := NewAdminUser(cfg, &fakes.FakeCmdStarter{})
			Expect(user.Username()).To(Equal("admin"))
			Expect(user.Password()).To(Equal("admin-password"))

			cfg = &config.Config{AdminUser: "admin-user-2", AdminPassword: "admin-password-2"}
			user = NewAdminUser(cfg, &fakes.FakeCmdStarter{})
			Expect(user.Username()).To(Equal("admin-user-2"))
			Expect(user.Password()).To(Equal("admin-password-2"))
		})
	})

	Describe("CreateUser", func() {
		var user *TestUser
		var fakeStarter *fakes.FakeCmdStarter
		var timeoutScale float64

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			timeoutScale = 1.0
		})

		JustBeforeEach(func() {
			cfg = &config.Config{
				TimeoutScale:         timeoutScale,
				UseExistingUser:      true,
				ExistingUser:         "my-username",
				ExistingUserPassword: "my-password",
			}

			user = NewTestUser(cfg, fakeStarter)
		})

		It("creates the user", func() {
			user.Create()
			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"create-user", user.Username(), user.Password()}))
		})

		Context("when 'cf create-user' exits with a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					user.Create()
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("(?s)Failed to create user.*scim_resource_already_exists"))
			})

			Context("and the output mentions that the user already exists", func() {
				BeforeEach(func() {
					fakeStarter.ToReturn[0].Output = "scim_resource_already_exists"
				})

				It("considers the command successful and does not fail", func() {
					failures := InterceptGomegaFailures(func() {
						user.Create()
					})
					Expect(failures).To(BeEmpty())
				})
			})

			Context("and it redacts the password", func() {
				JustBeforeEach(func() {
					fakeStarter.ToReturn[0].Output = fmt.Sprintf("blah blah %s %s", cfg.ExistingUser, cfg.ExistingUserPassword)
				})

				It("redactos", func() {
					failures := InterceptGomegaFailures(func() {
						user.Create()
					})
					Expect(failures[0]).NotTo(ContainSubstring(cfg.ExistingUserPassword))
					Expect(failures[0]).To(ContainSubstring("[REDACTED]"))
				})
			})

			Context("and stderr mentions that the user already exists", func() {
				BeforeEach(func() {
					fakeStarter.ToReturn[0].Stderr = "scim_resource_already_exists"
				})

				It("considers the command successful and does not fail", func() {
					failures := InterceptGomegaFailures(func() {
						user.Create()
					})
					Expect(failures).To(BeEmpty())
				})
			})
		})

		Context("when 'cf create-user' takes longer than the short timeout", func() {
			BeforeEach(func() {
				timeoutScale = 0.0334 // two-second timeout
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					user.Create()
				})

				Expect(len(failures)).To(BeNumerically(">", 0))
				Expect(failures[0]).To(MatchRegexp("(?s)Timed out after 2.*Failed to create user"))
			})
		})
	})

	Describe("Destroy", func() {
		var user *TestUser
		var fakeStarter *fakes.FakeCmdStarter
		var timeoutScale float64

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			timeoutScale = 1.0
		})

		JustBeforeEach(func() {
			cfg = &config.Config{
				TimeoutScale: timeoutScale,
			}
			user = NewTestUser(cfg, fakeStarter)
		})

		It("deletes the user", func() {
			user.Destroy()
			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"delete-user", "-f", user.Username()}))
		})

		Context("when 'cf delete-user' exits with a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					user.Destroy()
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("(?s)Failed to delete user.*to match exit code:.*0"))
			})
		})

		Context("when 'cf delete-user' times out", func() {
			BeforeEach(func() {
				timeoutScale = 0.0334 // two second timeout
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					user.Destroy()
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("(?s)Timed out after 2.*Failed to delete user"))
			})
		})
	})

	Describe("ShouldRemain", func() {
		var user *TestUser
		var fakeStarter *fakes.FakeCmdStarter
		var timeoutScale float64
		var shouldKeepUser bool

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			timeoutScale = 1.0
			shouldKeepUser = false
		})

		JustBeforeEach(func() {
			cfg = &config.Config{
				TimeoutScale:   timeoutScale,
				ShouldKeepUser: shouldKeepUser,
			}
			user = NewTestUser(cfg, fakeStarter)
		})

		It("returns false", func() {
			Expect(user.ShouldRemain()).To(BeFalse())
		})

		Context("when the config specifies that the user should not be deleted", func() {
			BeforeEach(func() {
				shouldKeepUser = true
			})

			It("returns true", func() {
				Expect(user.ShouldRemain()).To(BeTrue())
			})
		})
	})
})

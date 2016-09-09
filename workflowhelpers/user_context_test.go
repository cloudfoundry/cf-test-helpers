package workflowhelpers_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/cf-test-helpers/internal/fakes"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
)

var _ = Describe("UserContext", func() {
	Describe("NewUserContext", func() {
		var createUser = func() workflowhelpers.UserContext {
			return workflowhelpers.NewUserContext("http://FAKE_API.example.com", "FAKE_USERNAME", "FAKE_PASSWORD", "FAKE_ORG", "FAKE_SPACE", false)
		}

		It("returns a UserContext struct", func() {
			Expect(createUser()).To(BeAssignableToTypeOf(workflowhelpers.UserContext{}))
		})

		It("sets UserContext.ApiUrl", func() {
			Expect(createUser().ApiUrl).To(Equal("http://FAKE_API.example.com"))
		})

		It("sets UserContext.Username", func() {
			Expect(createUser().Username).To(Equal("FAKE_USERNAME"))
		})

		It("sets UserContext.Password", func() {
			Expect(createUser().Password).To(Equal("FAKE_PASSWORD"))
		})

		It("sets UserContext.Org", func() {
			Expect(createUser().Org).To(Equal("FAKE_ORG"))
		})

		It("sets UserContext.Space", func() {
			Expect(createUser().Space).To(Equal("FAKE_SPACE"))
		})
	})

	Describe("Login", func() {
		var target, username, password, org, space string
		var skipSslValidation bool
		var timeout time.Duration
		var fakeStarter *fakes.FakeCmdStarter

		var userContext workflowhelpers.UserContext

		BeforeEach(func() {
			target = "http://FAKE_API.example.com"
			username = "FAKE_USERNAME"
			password = "FAKE_PASSWORD"
			org = "FAKE_ORG"
			space = "FAKE_SPACE"
			skipSslValidation = false
			timeout = 1 * time.Second
			fakeStarter = fakes.NewFakeCmdStarter()
		})

		JustBeforeEach(func() {
			userContext = workflowhelpers.NewUserContext(target, username, password, org, space, skipSslValidation)
			userContext.CommandStarter = fakeStarter
		})

		It("logs in the user", func() {
			userContext.Login(timeout)

			Expect(fakeStarter.CalledWith).To(HaveLen(2))

			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"api", target}))

			Expect(fakeStarter.CalledWith[1].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[1].Args).To(Equal([]string{"auth", username, password}))
		})

		Context("when SkipSSLValidation is true", func() {
			BeforeEach(func() {
				skipSslValidation = true
			})

			It("adds the --skip-ssl-validation flag to 'cf api'", func() {
				userContext.Login(timeout)

				Expect(fakeStarter.CalledWith).To(HaveLen(2))

				Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
				Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"api", target, "--skip-ssl-validation"}))
			})
		})

		Context("when the 'cf api' call fails", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Login(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
			})
		})

		Context("when 'cf api' times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Login(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})

		})

		Context("when the 'cf auth' call fails", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[1].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Login(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
			})
		})

		Context("when the 'cf auth' times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[1].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Login(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})
		})
	})

	Describe("SetCfHomeDir", func() {
		var userContext workflowhelpers.UserContext
		var previousCfHome string
		BeforeEach(func() {
			previousCfHome = "my-cf-home-dir"
			os.Setenv("CF_HOME", previousCfHome)

			userContext = workflowhelpers.UserContext{}
		})

		AfterEach(func() {
			os.Unsetenv("CF_HOME")
		})

		It("creates a temporary directory and sets CF_HOME to point to it", func() {
			tmpDirRegexp := fmt.Sprintf("\\/var\\/folders\\/.*\\/.*\\/T\\/cf_home_%d", ginkgoconfig.GinkgoConfig.ParallelNode)

			userContext.SetCfHomeDir()
			cfHome := os.Getenv("CF_HOME")
			Expect(cfHome).To(MatchRegexp(tmpDirRegexp))
			Expect(cfHome).To(BeADirectory())
		})

		It("returns both the original and currently-used cf home directory", func() {
			originalCfHomeDir, currentCfHomeDir := userContext.SetCfHomeDir()
			Expect(originalCfHomeDir).To(Equal(previousCfHome))
			Expect(currentCfHomeDir).To(Equal(os.Getenv("CF_HOME")))
		})

		It("sets a unique CF_HOME value", func() {
			var firstHome, secondHome string

			_, firstHome = userContext.SetCfHomeDir()
			_, secondHome = userContext.SetCfHomeDir()

			Expect(firstHome).NotTo(Equal(secondHome))
		})

	})

	Describe("TargetSpace", func() {
		var userContext workflowhelpers.UserContext
		var org, space string
		var timeout time.Duration
		var fakeStarter *fakes.FakeCmdStarter

		BeforeEach(func() {
			org = "my-org"
			space = "my-space"
			timeout = 1 * time.Second
			fakeStarter = fakes.NewFakeCmdStarter()
		})

		JustBeforeEach(func() {
			userContext = workflowhelpers.UserContext{
				Org:            org,
				Space:          space,
				CommandStarter: fakeStarter,
			}
		})

		It("targets the org and space", func() {
			userContext.TargetSpace(timeout)
			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"target", "-o", org, "-s", space}))

			userContext.Org = "my-other-org"
			userContext.Space = "my-other-space"
			userContext.TargetSpace(timeout)
			Expect(fakeStarter.CalledWith[1].Args).To(Equal([]string{"target", "-o", "my-other-org", "-s", "my-other-space"}))
		})

		Context("when the space is not set", func() {
			BeforeEach(func() {
				space = ""
			})

			It("targets only the org", func() {
				userContext.TargetSpace(timeout)
				Expect(fakeStarter.CalledWith).To(HaveLen(1))
				Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
				Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"target", "-o", org}))
			})
		})

		Context("when the org is not set", func() {
			BeforeEach(func() {
				org = ""
			})

			It("does nothing", func() {
				userContext.TargetSpace(timeout)
				Expect(fakeStarter.CalledWith).To(HaveLen(0))
			})
		})

		Context("when the target command times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.TargetSpace(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})
		})

		Context("when the target command returns a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.TargetSpace(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
			})
		})
	})

	Describe("Logout", func() {
		var userContext workflowhelpers.UserContext
		var fakeStarter *fakes.FakeCmdStarter
		var timeout time.Duration

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			userContext = workflowhelpers.NewUserContext("", "", "", "", "", false)
			userContext.CommandStarter = fakeStarter
			timeout = 1 * time.Second
		})

		It("logs out the user", func() {
			userContext.Logout(timeout)

			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"logout"}))
		})

		Context("when 'cf logout' exits with a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Logout(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
			})
		})

		Context("when 'cf logout' times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.Logout(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})
		})
	})

	Describe("UnsetCfHomeDir", func() {
		var userContext workflowhelpers.UserContext
		var originalCfHomeDir, currentCfHomeDir string

		BeforeEach(func() {
			userContext = workflowhelpers.NewUserContext("", "", "", "", "", false)
		})

		It("restores Cf home dir to its original value", func() {
			originalCfHomeDir, currentCfHomeDir = userContext.SetCfHomeDir()
			userContext.UnsetCfHomeDir(originalCfHomeDir, currentCfHomeDir)
			Expect(os.Getenv("CF_HOME")).To(Equal(originalCfHomeDir))
			Expect(currentCfHomeDir).NotTo(BeADirectory())
		})
	})

	Describe("Destroy", func() {
		var userContext workflowhelpers.UserContext
		var fakeStarter *fakes.FakeCmdStarter
		var timeout time.Duration

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			userContext = workflowhelpers.NewUserContext("", "", "", "", "", false)
			userContext.CommandStarter = fakeStarter
			timeout = 1 * time.Second
		})

		It("deletes the user", func() {
			userContext.DeleteUser(timeout)
			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"delete-user", "-f", userContext.Username}))
		})

		Context("when 'cf delete-user' exits with a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.DeleteUser(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
			})
		})

		Context("when 'cf delete-user' times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.DeleteUser(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})
		})

		XContext("when config.ShouldKeepUser is true", func() {
			JustBeforeEach(func() {

			})

			It("does not delete the user", func() {

			})
		})
	})

	Describe("create-user", func() {
		var userContext workflowhelpers.UserContext
		var fakeStarter *fakes.FakeCmdStarter
		var timeout time.Duration

		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
			userContext = workflowhelpers.NewUserContext("", "", "", "", "", false)
			userContext.CommandStarter = fakeStarter
			timeout = 1 * time.Second
		})

		It("creates the user", func() {
			userContext.CreateUser(timeout)
			Expect(fakeStarter.CalledWith).To(HaveLen(1))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"create-user", userContext.Username, userContext.Password}))
		})

		Context("when 'cf create-user' exits with a non-zero exit code", func() {
			BeforeEach(func() {
				fakeStarter.ToReturn[0].ExitCode = 1
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.CreateUser(timeout)
				})

				Expect(failures).To(HaveLen(1))
				Expect(failures[0]).To(MatchRegexp("scim_resource_already_exists"))
			})

			Context("and the output mentions that the user already exists", func() {
				BeforeEach(func() {
					fakeStarter.ToReturn[0].Output = "scim_resource_already_exists"
				})

				It("considers the command successful and does not fail", func() {
					failures := InterceptGomegaFailures(func() {
						userContext.CreateUser(timeout)
					})
					Expect(failures).To(BeEmpty())
				})
			})
		})

		Context("when 'cf create-user' times out", func() {
			BeforeEach(func() {
				timeout = 2 * time.Second
				fakeStarter.ToReturn[0].SleepTime = 3
			})

			It("fails with a ginkgo error", func() {
				failures := InterceptGomegaFailures(func() {
					userContext.CreateUser(timeout)
				})

				Expect(len(failures)).To(BeNumerically(">", 0))
				Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
			})
		})

		XContext("when config.UseExistingUser is true", func() {
			JustBeforeEach(func() {

			})

			It("does not delete the user", func() {

			})
		})
	})
})

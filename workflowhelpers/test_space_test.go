package workflowhelpers_test

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal/fakes"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TestSpace", func() {
	var cfg config.Config
	var userContext UserContext
	var namePrefix string
	BeforeEach(func() {
		namePrefix = "UNIT-TEST"
		cfg = config.Config{
			NamePrefix:   namePrefix,
			TimeoutScale: 1.0,
		}
		userContext = NewUserContext("api.url", "my-user", "my-password", "my-org", "my-space", true)
	})

	Describe("NewRegularTestSpace", func() {
		It("generates a quotaDefinitionName", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionName).To(MatchRegexp("%s-[0-9]-QUOTA-.*", namePrefix))
		})

		It("generates an organizationName", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.OrganizationName).To(MatchRegexp("%s-[0-9]-ORG-.*", namePrefix))
		})

		It("generates a spaceName", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.SpaceName).To(MatchRegexp("%s-[0-9]-SPACE-.*", namePrefix))
		})

		It("sets a timeout for cf commands", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.Timeout).To(Equal(1 * time.Minute))
		})

		It("sets the UserContext", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.UserContext).To(Equal(userContext))
		})

		Context("when the config scales the timeout", func() {
			BeforeEach(func() {
				cfg = config.Config{
					NamePrefix:   namePrefix,
					TimeoutScale: 2.0,
				}
			})

			It("scales the timeout for cf commands", func() {
				testSpace := NewRegularTestSpace(cfg)
				Expect(testSpace.Timeout).To(Equal(2 * time.Minute))
			})
		})

		It("uses default values for the quota", func() {
			testSpace := NewRegularTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionTotalMemoryLimit).To(Equal("10G"))
			Expect(testSpace.QuotaDefinitionInstanceMemoryLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionRoutesLimit).To(Equal("1000"))
			Expect(testSpace.QuotaDefinitionAppInstanceLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionServiceInstanceLimit).To(Equal("100"))
			Expect(testSpace.QuotaDefinitionAllowPaidServicesFlag).To(Equal("--allow-paid-service-plans"))
		})
	})

	Describe("NewPersistentAppTestSpace", func() {
		var quotaDefinitionName, organizationName, spaceName string
		BeforeEach(func() {
			quotaDefinitionName = "persistent-quota"
			organizationName = "persistent-org"
			spaceName = "persistent-space"
			cfg = config.Config{
				PersistentAppOrg:       organizationName,
				PersistentAppSpace:     spaceName,
				PersistentAppQuotaName: quotaDefinitionName,
			}
		})

		It("gets the quota definition name from the config", func() {
			testSpace := NewPersistentAppTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionName).To(Equal(quotaDefinitionName))
		})

		It("gets the org name from the config", func() {
			testSpace := NewPersistentAppTestSpace(cfg)
			Expect(testSpace.OrganizationName).To(Equal(organizationName))
		})

		It("gets the space name from the config", func() {
			testSpace := NewPersistentAppTestSpace(cfg)
			Expect(testSpace.SpaceName).To(Equal(spaceName))
		})

		It("uses default values for the quota", func() {
			testSpace := NewPersistentAppTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionTotalMemoryLimit).To(Equal("10G"))
			Expect(testSpace.QuotaDefinitionInstanceMemoryLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionRoutesLimit).To(Equal("1000"))
			Expect(testSpace.QuotaDefinitionAppInstanceLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionServiceInstanceLimit).To(Equal("100"))
			Expect(testSpace.QuotaDefinitionAllowPaidServicesFlag).To(Equal("--allow-paid-service-plans"))
		})

		It("sets isPersistent to be true", func() {
			testSpace := NewPersistentAppTestSpace(cfg)
			Expect(testSpace.IsPersistent).To(Equal(true))
		})
	})

	Describe("NewRunawayAppTestSpace", func() {
		BeforeEach(func() {
			cfg = config.Config{
				NamePrefix: "UNIT-TEST",
			}
		})
		It("generates an org name, space name, and quota name", func() {
			testSpace := NewRunawayAppTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionName).To(MatchRegexp("%s-[0-9]-QUOTA-.*", namePrefix))
			Expect(testSpace.OrganizationName).To(MatchRegexp("%s-[0-9]-ORG-.*", namePrefix))
			Expect(testSpace.SpaceName).To(MatchRegexp("%s-[0-9]-SPACE-.*", namePrefix))
		})

		It("sets the quota definition memory limit to a very high number", func() {
			testSpace := NewRunawayAppTestSpace(cfg)
			Expect(testSpace.QuotaDefinitionTotalMemoryLimit).To(Equal(RUNAWAY_QUOTA_MEM_LIMIT))

			Expect(testSpace.QuotaDefinitionInstanceMemoryLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionRoutesLimit).To(Equal("1000"))
			Expect(testSpace.QuotaDefinitionAppInstanceLimit).To(Equal("-1"))
			Expect(testSpace.QuotaDefinitionServiceInstanceLimit).To(Equal("100"))
			Expect(testSpace.QuotaDefinitionAllowPaidServicesFlag).To(Equal("--allow-paid-service-plans"))
		})
	})

	Describe("InstantiateRemotely", func() {
		var testSpace *TestSpace
		var fakeStarter *fakes.FakeCmdStarter
		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
		})

		JustBeforeEach(func() {
			testSpace = NewRegularTestSpace(cfg)
			testSpace.CommandStarter = fakeStarter
		})

		It("creates a quota", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 0))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{
				"create-quota", testSpace.QuotaDefinitionName,
				"-m", testSpace.QuotaDefinitionTotalMemoryLimit,
				"-i", testSpace.QuotaDefinitionInstanceMemoryLimit,
				"-r", testSpace.QuotaDefinitionRoutesLimit,
				"-a", testSpace.QuotaDefinitionAppInstanceLimit,
				"-s", testSpace.QuotaDefinitionServiceInstanceLimit,
				testSpace.QuotaDefinitionAllowPaidServicesFlag,
			}))
		})

		It("creates an org", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 1))
			Expect(fakeStarter.CalledWith[1].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[1].Args).To(Equal([]string{"create-org", testSpace.OrganizationName}))
		})

		It("sets quota", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 2))
			Expect(fakeStarter.CalledWith[2].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[2].Args).To(Equal([]string{"set-quota", testSpace.OrganizationName, testSpace.QuotaDefinitionName}))
		})

		It("create space", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 3))
			Expect(fakeStarter.CalledWith[3].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[3].Args).To(Equal([]string{"create-space", "-o", testSpace.OrganizationName, testSpace.SpaceName}))
		})

		It("sets space manager", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 4))
			Expect(fakeStarter.CalledWith[4].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[4].Args).To(Equal([]string{"set-space-role", testSpace.RegularUserUsername, testSpace.OrganizationName, testSpace.SpaceName, "SpaceManager"}))
		})

		It("sets space developer", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 5))
			Expect(fakeStarter.CalledWith[5].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[5].Args).To(Equal([]string{"set-space-role", testSpace.RegularUserUsername, testSpace.OrganizationName, testSpace.SpaceName, "SpaceDeveloper"}))
		})

		It("sets space auditor", func() {
			testSpace.InstantiateRemotely()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 6))
			Expect(fakeStarter.CalledWith[6].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[6].Args).To(Equal([]string{"set-space-role", testSpace.RegularUserUsername, testSpace.OrganizationName, testSpace.SpaceName, "SpaceAuditor"}))
		})

		Describe("failure cases", func() {
			testFailureCase := func(callIndex int) func() {
				return func() {
					BeforeEach(func() {
						fakeStarter.ToReturn[callIndex].ExitCode = 1
					})

					It("returns a ginkgo error", func() {
						failures := InterceptGomegaFailures(func() {
							testSpace.InstantiateRemotely()
						})
						Expect(failures).To(HaveLen(1))
						Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
					})
				}
			}

			Context("when 'cf create-quota' fails", testFailureCase(0))
			Context("when 'cf create-org' fails", testFailureCase(1))
			Context("when 'cf set-quota' fails", testFailureCase(2))
			Context("when 'cf create-space' fails", testFailureCase(3))
			Context("when 'cf set-space-manager' fails", testFailureCase(4))
			Context("when 'cf set-space-developer", testFailureCase(5))
			Context("when 'cf set-space-auditor' fails", testFailureCase(6))
		})

		Describe("timing out", func() {
			BeforeEach(func() {
				cfg = config.Config{
					TimeoutScale: 0.03333, // 2 second timeout
				}
			})

			testTimeoutCase := func(callIndex int) func() {
				return func() {
					BeforeEach(func() {
						fakeStarter.ToReturn[callIndex].SleepTime = 5
					})

					It("returns a ginkgo error", func() {
						failures := InterceptGomegaFailures(func() {
							testSpace.InstantiateRemotely()
						})

						Expect(failures).To(HaveLen(1))
						Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
					})
				}
			}

			Context("when 'cf create-quota' times out", testTimeoutCase(0))
			Context("when 'cf create-org' times out", testTimeoutCase(1))
			Context("when 'cf set-quota' times out", testTimeoutCase(2))
			Context("when 'cf create-space' times out", testTimeoutCase(3))
			Context("when 'cf set-space-manager' times out", testTimeoutCase(4))
			Context("when 'cf set-space-developer", testTimeoutCase(5))
			Context("when 'cf set-space-auditor' times out", testTimeoutCase(6))
		})

	})

	Describe("Destroy", func() {
		var testSpace *TestSpace
		var fakeStarter *fakes.FakeCmdStarter
		BeforeEach(func() {
			fakeStarter = fakes.NewFakeCmdStarter()
		})

		JustBeforeEach(func() {
			testSpace = NewRegularTestSpace(cfg)
			testSpace.CommandStarter = fakeStarter
		})

		It("deletes the org", func() {
			testSpace.Destroy()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 0))
			Expect(fakeStarter.CalledWith[0].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[0].Args).To(Equal([]string{"delete-org", "-f", testSpace.OrganizationName}))
		})

		It("deletes the space", func() {
			testSpace.Destroy()
			Expect(len(fakeStarter.CalledWith)).To(BeNumerically(">", 1))
			Expect(fakeStarter.CalledWith[1].Executable).To(Equal("cf"))
			Expect(fakeStarter.CalledWith[1].Args).To(Equal([]string{"delete-quota", "-f", testSpace.QuotaDefinitionName}))
		})

		Context("when context.isPersistent is true", func() {
			BeforeEach(func() {
				testSpace.IsPersistent = true
			})

			It("does not delete the org", func() {
				Expect(len(fakeStarter.CalledWith)).To(Equal(0))
			})

			It("does not delete the space", func() {
				Expect(len(fakeStarter.CalledWith)).To(Equal(0))
			})
		})

		Describe("failure cases", func() {
			testFailureCase := func(callIndex int) func() {
				return func() {
					BeforeEach(func() {
						fakeStarter.ToReturn[callIndex].ExitCode = 1
					})

					It("returns a ginkgo error", func() {
						failures := InterceptGomegaFailures(func() {
							testSpace.Destroy()
						})
						Expect(failures).To(HaveLen(1))
						Expect(failures[0]).To(MatchRegexp("to match exit code:\n.*0"))
					})
				}
			}

			Context("when 'delete-org' fails", testFailureCase(0))
			Context("when 'delete-quota' fails", testFailureCase(1))
		})

		Describe("timing out", func() {
			BeforeEach(func() {
				cfg = config.Config{
					TimeoutScale: 0.03333, // 2 second timeout
				}
			})

			testTimeoutCase := func(callIndex int) func() {
				return func() {
					BeforeEach(func() {
						fakeStarter.ToReturn[callIndex].SleepTime = 5
					})

					It("returns a ginkgo error", func() {
						failures := InterceptGomegaFailures(func() {
							testSpace.Destroy()
						})

						Expect(failures).To(HaveLen(1))
						Expect(failures[0]).To(MatchRegexp("Timed out after 2.*"))
					})
				}
			}

			Context("when 'cf delete-org' times out", testTimeoutCase(0))
			Context("when 'cf delete-quota' times out", testTimeoutCase(1))
		})
	})
})

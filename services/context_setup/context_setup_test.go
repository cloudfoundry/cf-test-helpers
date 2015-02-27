package context_setup_test

import (
	"fmt"
	"os/exec"
	"encoding/json"

	. "github.com/cloudfoundry-incubator/cf-test-helpers/services/context_setup"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
)

type apiRequestInputs struct {
	method, endpoint string
	data []string
}

var _ = Describe("ConfiguredContext", func() {
	var (
		FakeCfCalls [][]string
		FakeCfCallback map[string]func(args ...string) *exec.Cmd
		FakeAsUserCalls []cf.UserContext
		FakeApiRequestCalls []apiRequestInputs
		FakeApiRequestCallbacks []func(method, endpoint string, response interface{}, data ...string)
	)

	var FakeCf = func(args ...string) *gexec.Session {
		FakeCfCalls = append(FakeCfCalls, args)
		var cmd *exec.Cmd
		if callback, exists := FakeCfCallback[args[0]]; exists {
			cmd = callback(args...)
		} else {
			cmd = exec.Command("echo", "OK")
		}
		session, _ := gexec.Start(cmd, nil, nil)
		return session
	}

	var FakeAsUser = func(userContext cf.UserContext, actions func()) {
		FakeAsUserCalls = append(FakeAsUserCalls, userContext)
		actions()
	}

	var FakeApiRequest = func(method, endpoint string, response interface{}, data ...string) {
		FakeApiRequestCalls = append(FakeApiRequestCalls, apiRequestInputs{
			method: method,
			endpoint: endpoint,
			data: data,
		})

		if len(FakeApiRequestCallbacks) > 0 {
			callback := FakeApiRequestCallbacks[0]
			FakeApiRequestCallbacks = FakeApiRequestCallbacks[1:]
			callback(method, endpoint, response, data...)
		}
	}

	BeforeEach(func() {
		FakeCfCalls = [][]string{}
		FakeCfCallback = map[string]func(args ...string) *exec.Cmd{}
		FakeAsUserCalls = []cf.UserContext{}
		FakeApiRequestCalls = []apiRequestInputs{}
		FakeApiRequestCallbacks = []func(method, endpoint string, response interface{}, data ...string){}

		cf.Cf = FakeCf
		cf.AsUser = FakeAsUser
		cf.ApiRequest = FakeApiRequest

		TimeoutScale = float64(0.1)
	})

	Describe("Setup", func() {
		var (
			config IntegrationConfig
			prefix = "fake-prefix"

			context *ConfiguredContext
		)

		BeforeEach(func() {
			config = IntegrationConfig{
				AppsDomain: "fake-domain",
				ApiEndpoint: "fake-endpoint",
				AdminUser: "fake-admin-user",
				AdminPassword: "fake-admin-password",
				SkipSSLValidation: true,
			}
			context = NewContext(config, prefix)
		})

		It("executes commands as the admin user", func() {
			context.Setup()

			Expect(FakeAsUserCalls).To(Equal([]cf.UserContext{
				{
					ApiUrl: "fake-endpoint",
					Username: "fake-admin-user",
					Password: "fake-admin-password",
					SkipSSLValidation: true,
				},
			}))
		})

		It("creates a new user with a unique name", func() {
			context.Setup()

			createUserCall := FakeCfCalls[0]
			Expect(createUserCall[0]).To(Equal("create-user"))
			Expect(createUserCall[1]).To(MatchRegexp("fake-prefix-USER-\\d+-.*"))
			Expect(createUserCall[2]).To(Equal("meow")) //why meow??
		})

		It("creates a new quota with a unique name", func() {
			FakeApiRequestCallbacks = append(FakeApiRequestCallbacks, func(method, endpoint string, response interface{}, data ...string){
				genericResponse, ok := response.(*cf.GenericResource)
				if !ok {
					Fail(fmt.Sprintf("Expected response to be of type *cf.GenericResource: %#v", response))
				}
				genericResponse.Metadata.Guid = "fake-guid"
			})

			context.Setup()

			createQuotaCall := FakeApiRequestCalls[0]
			Expect(createQuotaCall.method).To(Equal("POST"))
			Expect(createQuotaCall.endpoint).To(Equal("/v2/quota_definitions"))
			Expect(createQuotaCall.data).To(HaveLen(1))

			definitionJSON := createQuotaCall.data[0]
			definition := &QuotaDefinition{}
			err := json.Unmarshal([]byte(definitionJSON), definition)
			Expect(err).ToNot(HaveOccurred())

			Expect(definition.Name).To(MatchRegexp("fake-prefix-QUOTA-\\d+-.*"))
			Expect(definition.TotalServices).To(Equal(100))
			Expect(definition.TotalRoutes).To(Equal(1000))
			Expect(definition.MemoryLimit).To(Equal(10240))
			Expect(definition.NonBasicServicesAllowed).To(BeTrue())
		})

		It("creates a new org with a unique name", func() {
			context.Setup()

			createOrgCall := FakeCfCalls[1]
			Expect(createOrgCall).To(HaveLen(2))
			Expect(createOrgCall[0]).To(Equal("create-org"))
			Expect(createOrgCall[1]).To(MatchRegexp("fake-prefix-ORG-\\d+-.*"))
		})

		It("set the quota on the new org", func() {
			context.Setup()

			setQuotaCall := FakeCfCalls[2]
			Expect(setQuotaCall).To(HaveLen(3))
			Expect(setQuotaCall[0]).To(Equal("set-quota"))
			Expect(setQuotaCall[1]).To(MatchRegexp("fake-prefix-ORG-\\d+-.*"))
			Expect(setQuotaCall[2]).To(MatchRegexp("fake-prefix-QUOTA-\\d+-.*"))
		})

		Context("when create-user fails", func() {
			BeforeEach(func() {
				FakeCfCallback["create-user"] = func(args ...string) *exec.Cmd {
					return exec.Command("false")
				}
			})

			It("causes gomega matcher failure with stdout & stderr", func() {
				failures := InterceptGomegaFailures(func() {
					context.Setup()
				})
				Expect(failures[0]).To(MatchRegexp("Failed to create user \\(exit 1\\):\n\n\\[stdout\\]:\n.*\n\n\\[stderr\\]:\n.*"))
			})
		})

		Context("when create-org fails", func() {
			BeforeEach(func() {
				FakeCfCallback["create-org"] = func(args ...string) *exec.Cmd {
					return exec.Command("false")
				}
			})

			It("causes gomega matcher failure with stdout & stderr", func() {
				failures := InterceptGomegaFailures(func() {
					context.Setup()
				})
				Expect(failures[0]).To(MatchRegexp("Failed to create org \\(exit 1\\):\n\n\\[stdout\\]:\n.*\n\n\\[stderr\\]:\n.*"))
			})
		})

		Context("when set-quota fails", func() {
			BeforeEach(func() {
				FakeCfCallback["set-quota"] = func(args ...string) *exec.Cmd {
					return exec.Command("false")
				}
			})

			It("causes gomega matcher failure with stdout & stderr", func() {
				failures := InterceptGomegaFailures(func() {
					context.Setup()
				})
				Expect(failures[0]).To(MatchRegexp("Failed to set org quota \\(exit 1\\):\n\n\\[stdout\\]:\n.*\n\n\\[stderr\\]:\n.*"))
			})
		})
	})

	Describe("Teardown", func() {
		var (
			config IntegrationConfig
			prefix = "fake-prefix"

			context *ConfiguredContext
		)

		BeforeEach(func() {
			config = IntegrationConfig{
				AppsDomain: "fake-domain",
				ApiEndpoint: "fake-endpoint",
				AdminUser: "fake-admin-user",
				AdminPassword: "fake-admin-password",
				SkipSSLValidation: true,
			}
			context = NewContext(config, prefix)

			FakeApiRequestCallbacks = append(FakeApiRequestCallbacks, func(method, endpoint string, response interface{}, data ...string){
				genericResponse, ok := response.(*cf.GenericResource)
				if !ok {
					Fail(fmt.Sprintf("Expected response to be of type *cf.GenericResource: %#v", response))
				}
				genericResponse.Metadata.Guid = "fake-guid"
			})

			context.Setup()

			// ignore calls made by setup
			FakeCfCalls = [][]string{}
			FakeAsUserCalls = []cf.UserContext{}
			FakeApiRequestCalls = []apiRequestInputs{}
		})

		It("executes commands as the admin user", func() {
			context.Teardown()

			Expect(FakeAsUserCalls).To(Equal([]cf.UserContext{
				{
					ApiUrl: "fake-endpoint",
					Username: "fake-admin-user",
					Password: "fake-admin-password",
					SkipSSLValidation: true,
				},
			}))
		})

		It("deletes the user created by Setup()", func() {
			context.Teardown()

			deleteUserCall := FakeCfCalls[0]
			Expect(deleteUserCall[0]).To(Equal("delete-user"))
			Expect(deleteUserCall[1]).To(Equal("-f"))
			Expect(deleteUserCall[2]).To(MatchRegexp("fake-prefix-USER-\\d+-.*"))
		})

		It("deletes the org created by Setup()", func() {
			context.Teardown()

			deleteOrgCall := FakeCfCalls[1]
			Expect(deleteOrgCall[0]).To(Equal("delete-org"))
			Expect(deleteOrgCall[1]).To(Equal("-f"))
			Expect(deleteOrgCall[2]).To(MatchRegexp("fake-prefix-ORG-\\d+-.*"))
		})

		It("deletes the quota created by Setup()", func() {
			context.Teardown()

			deleteQuotaCall := FakeApiRequestCalls[0]
			Expect(deleteQuotaCall.method).To(Equal("DELETE"))
			Expect(deleteQuotaCall.endpoint).To(Equal("/v2/quota_definitions/fake-guid?recursive=true"))
		})

		Context("when delete-user fails", func() {
			BeforeEach(func() {
				FakeCfCallback["delete-user"] = func(args ...string) *exec.Cmd {
					return exec.Command("false")
				}
			})

			It("causes gomega matcher failure with stdout & stderr", func() {
				failures := InterceptGomegaFailures(func() {
					context.Teardown()
				})
				Expect(failures[0]).To(MatchRegexp("Failed to delete user \\(exit 1\\):\n\n\\[stdout\\]:\n.*\n\n\\[stderr\\]:\n.*"))
			})
		})

		Context("when delete-org fails", func() {
			BeforeEach(func() {
				FakeCfCallback["delete-org"] = func(args ...string) *exec.Cmd {
					return exec.Command("false")
				}
			})

			It("causes gomega matcher failure with stdout & stderr", func() {
				failures := InterceptGomegaFailures(func() {
					context.Teardown()
				})
				Expect(failures[0]).To(MatchRegexp("Failed to delete org \\(exit 1\\):\n\n\\[stdout\\]:\n.*\n\n\\[stderr\\]:\n.*"))
			})
		})
	})
})

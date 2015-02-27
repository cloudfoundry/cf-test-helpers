package context_setup

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

type ConfiguredContext struct {
	config IntegrationConfig

	organizationName string
	spaceName        string

	quotaDefinitionName string
	quotaDefinitionGUID string

	regularUserUsername string
	regularUserPassword string

	isPersistent bool
}

type QuotaDefinition struct {
	Name string `json:"name"`

	NonBasicServicesAllowed bool `json:"non_basic_services_allowed"`

	TotalServices int `json:"total_services"`
	TotalRoutes   int `json:"total_routes"`

	MemoryLimit int `json:"memory_limit"`
}

func NewContext(config IntegrationConfig, prefix string) *ConfiguredContext {
	node := ginkgoconfig.GinkgoConfig.ParallelNode
	timeTag := time.Now().Format("2006_01_02-15h04m05.999s")

	return &ConfiguredContext{
		config: config,

		quotaDefinitionName: fmt.Sprintf("%s-QUOTA-%d-%s", prefix, node, timeTag),

		organizationName: fmt.Sprintf("%s-ORG-%d-%s", prefix, node, timeTag),
		spaceName:        fmt.Sprintf("%s-SPACE-%d-%s", prefix, node, timeTag),

		regularUserUsername: fmt.Sprintf("%s-USER-%d-%s", prefix, node, timeTag),
		regularUserPassword: "meow",

		isPersistent: false,
	}
}

func (context *ConfiguredContext) Setup() {
	cf.AsUser(context.AdminUserContext(), func() {
		shortTimeout := ScaledTimeout(10*time.Second)
		createUserCmd := cf.Cf("create-user", context.regularUserUsername, context.regularUserPassword).Wait(shortTimeout)
		exitCode := createUserCmd.ExitCode()
		Expect(exitCode).ToNot(Equal(-1), "Timed out creating user (%s)", shortTimeout)
		Expect(exitCode).To(Equal(0), "Failed to create user (exit %d):\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, string(createUserCmd.Out.Contents()), string(createUserCmd.Err.Contents()))

		definition := QuotaDefinition{
			Name: context.quotaDefinitionName,

			TotalServices: 100,
			TotalRoutes:   1000,

			MemoryLimit: 10240,

			NonBasicServicesAllowed: true,
		}

		definitionPayload, err := json.Marshal(definition)
		Expect(err).ToNot(HaveOccurred())

		var response cf.GenericResource

		cf.ApiRequest("POST", "/v2/quota_definitions", &response, string(definitionPayload))

		context.quotaDefinitionGUID = response.Metadata.Guid

		longTimeout := ScaledTimeout(60*time.Second)

		createOrgCmd := cf.Cf("create-org", context.organizationName).Wait(longTimeout)
		exitCode = createOrgCmd.ExitCode()
		Expect(exitCode).ToNot(Equal(-1), "Timed out creating org (%s)", longTimeout)
		Expect(exitCode).To(Equal(0), "Failed to create org (exit %d):\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, string(createOrgCmd.Out.Contents()), string(createOrgCmd.Err.Contents()))

		setQuotaCmd := cf.Cf("set-quota", context.organizationName, context.quotaDefinitionName).Wait(longTimeout)
		exitCode = setQuotaCmd.ExitCode()
		Expect(exitCode).ToNot(Equal(-1), "Timed out setting org quota (%s)", longTimeout)
		Expect(exitCode).To(Equal(0), "Failed to set org quota (exit %d):\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, string(setQuotaCmd.Out.Contents()), string(setQuotaCmd.Err.Contents()))
	})
}

func (context *ConfiguredContext) Teardown() {
	cf.AsUser(context.AdminUserContext(), func() {
		longTimeout := ScaledTimeout(60*time.Second)
		deleteUserCmd := cf.Cf("delete-user", "-f", context.regularUserUsername).Wait(longTimeout)
		exitCode := deleteUserCmd.ExitCode()
		Expect(exitCode).ToNot(Equal(-1), "Timed out deleting user (%s)", longTimeout)
		Expect(exitCode).To(Equal(0), "Failed to delete user (exit %d):\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, string(deleteUserCmd.Out.Contents()), string(deleteUserCmd.Err.Contents()))

		if !context.isPersistent {
			deleteOrgCmd := cf.Cf("delete-org", "-f", context.organizationName).Wait(longTimeout)
			exitCode := deleteOrgCmd.ExitCode()
			Expect(exitCode).ToNot(Equal(-1), "Timed out deleting org (%s)", longTimeout)
			Expect(exitCode).To(Equal(0), "Failed to delete org (exit %d):\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, string(deleteOrgCmd.Out.Contents()), string(deleteOrgCmd.Err.Contents()))

			cf.ApiRequest(
				"DELETE",
				"/v2/quota_definitions/"+context.quotaDefinitionGUID+"?recursive=true",
				nil,
			)
		}
	})
}

func (context *ConfiguredContext) AdminUserContext() cf.UserContext {
	return cf.NewUserContext(
		context.config.ApiEndpoint,
		context.config.AdminUser,
		context.config.AdminPassword,
		"",
		"",
		context.config.SkipSSLValidation,
	)
}

func (context *ConfiguredContext) RegularUserContext() cf.UserContext {
	return cf.NewUserContext(
		context.config.ApiEndpoint,
		context.regularUserUsername,
		context.regularUserPassword,
		context.organizationName,
		context.spaceName,
		context.config.SkipSSLValidation,
	)
}

package context_setup

import (
	"encoding/json"
	"fmt"
	"time"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
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
		shortTimeout := ScaledTimeout(10 * time.Second)

		ExecWithTimeout(cf.Cf("create-user", context.regularUserUsername, context.regularUserPassword), shortTimeout)

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

		longTimeout := ScaledTimeout(60 * time.Second)

		ExecWithTimeout(cf.Cf("create-org", context.organizationName), longTimeout)
		ExecWithTimeout(cf.Cf("set-quota", context.organizationName, context.quotaDefinitionName), longTimeout)
	})
}

func (context *ConfiguredContext) Teardown() {
	cf.AsUser(context.AdminUserContext(), func() {
		longTimeout := ScaledTimeout(60 * time.Second)
		ExecWithTimeout(cf.Cf("delete-user", "-f", context.regularUserUsername), longTimeout)

		if !context.isPersistent {
			ExecWithTimeout(cf.Cf("delete-org", "-f", context.organizationName), longTimeout)

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

package services

import (
	"encoding/json"
	"fmt"
	"time"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

type Context interface {
    Setup()
    Teardown()

    AdminUserContext() cf.UserContext
    RegularUserContext() cf.UserContext

    ShortTimeout() time.Duration
    LongTimeout() time.Duration
}

type context struct {
	config Config

    shortTimeout time.Duration
    longTimeout time.Duration

	organizationName string
	spaceName        string

	quotaDefinitionName string
	quotaDefinitionGUID string

	regularUserUsername string
	regularUserPassword string

	isPersistent bool

    originalCfHomeDir string
    currentCfHomeDir string
}

type QuotaDefinition struct {
	Name string `json:"name"`

	NonBasicServicesAllowed bool `json:"non_basic_services_allowed"`

	TotalServices int `json:"total_services"`
	TotalRoutes   int `json:"total_routes"`

	MemoryLimit int `json:"memory_limit"`
}

func NewContext(config Config, prefix string) Context {
	node := ginkgoconfig.GinkgoConfig.ParallelNode
	timeTag := time.Now().Format("2006_01_02-15h04m05.999s")

	return &context{
		config: config,

        shortTimeout: config.ScaledTimeout(1 * time.Minute),
        longTimeout: config.ScaledTimeout(5 * time.Minute),

		quotaDefinitionName: fmt.Sprintf("%s-QUOTA-%d-%s", prefix, node, timeTag),

		organizationName: fmt.Sprintf("%s-ORG-%d-%s", prefix, node, timeTag),
		spaceName:        fmt.Sprintf("%s-SPACE-%d-%s", prefix, node, timeTag),

		regularUserUsername: fmt.Sprintf("%s-USER-%d-%s", prefix, node, timeTag),
		regularUserPassword: "meow",

		isPersistent: false,
	}
}

func (c context) ShortTimeout() time.Duration {
    return c.shortTimeout
}

func (c context) LongTimeout() time.Duration {
    return c.longTimeout
}

func (c *context) Setup() {
	cf.AsUser(c.AdminUserContext(), func() {
		runner.NewCmdRunner(cf.Cf("create-user", c.regularUserUsername, c.regularUserPassword), c.shortTimeout).Run()

		definition := QuotaDefinition{
			Name: c.quotaDefinitionName,

			TotalServices: 100,
			TotalRoutes:   1000,

			MemoryLimit: 10240,

			NonBasicServicesAllowed: true,
		}

		definitionPayload, err := json.Marshal(definition)
		Expect(err).ToNot(HaveOccurred())

		var response cf.GenericResource

		cf.ApiRequest("POST", "/v2/quota_definitions", &response, string(definitionPayload))

		c.quotaDefinitionGUID = response.Metadata.Guid

		runner.NewCmdRunner(cf.Cf("create-org", c.organizationName), c.shortTimeout).Run()
		runner.NewCmdRunner(cf.Cf("set-quota", c.organizationName, c.quotaDefinitionName), c.shortTimeout).Run()
	})

    cf.AsUser(c.AdminUserContext(), func() {
        c.setUpSpaceWithUserAccess(c.RegularUserContext())
    })

    c.originalCfHomeDir, c.currentCfHomeDir = cf.InitiateUserContext(c.RegularUserContext())
    cf.TargetSpace(c.RegularUserContext())
}

func (c *context) Teardown() {
    cf.RestoreUserContext(c.RegularUserContext(), c.originalCfHomeDir, c.currentCfHomeDir)

	cf.AsUser(c.AdminUserContext(), func() {
		runner.NewCmdRunner(cf.Cf("delete-user", "-f", c.regularUserUsername), c.longTimeout).Run()

		if !c.isPersistent {
			runner.NewCmdRunner(cf.Cf("delete-org", "-f", c.organizationName), c.longTimeout).Run()

			cf.ApiRequest(
				"DELETE",
				"/v2/quota_definitions/"+c.quotaDefinitionGUID+"?recursive=true",
				nil,
			)
		}
	})
}

func (c context) AdminUserContext() cf.UserContext {
	return cf.NewUserContext(
		c.config.ApiEndpoint,
		c.config.AdminUser,
		c.config.AdminPassword,
		"",
		"",
		c.config.SkipSSLValidation,
	)
}

func (c context) RegularUserContext() cf.UserContext {
	return cf.NewUserContext(
		c.config.ApiEndpoint,
		c.regularUserUsername,
		c.regularUserPassword,
		c.organizationName,
		c.spaceName,
		c.config.SkipSSLValidation,
	)
}

func (c context) setUpSpaceWithUserAccess(uc cf.UserContext) {
    runner.NewCmdRunner(cf.Cf("create-space", "-o", uc.Org, uc.Space), c.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceManager"), c.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceDeveloper"), c.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceAuditor"), c.shortTimeout).Run()
}

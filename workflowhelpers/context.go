package workflowhelpers

import (
	"fmt"
	"time"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
)

type ReproducibleTestSetup struct {
	config config.Config

	shortTimeout time.Duration
	longTimeout  time.Duration

	organizationName string
	spaceName        string

	quotaDefinitionName                  string
	quotaDefinitionTotalMemoryLimit      string
	quotaDefinitionInstanceMemoryLimit   string
	quotaDefinitionRoutesLimit           string
	quotaDefinitionAppInstanceLimit      string
	quotaDefinitionServiceInstanceLimit  string
	quotaDefinitionAllowPaidServicesFlag string

	regularUserUsername string
	regularUserPassword string

	adminUserUsername string
	adminUserPassword string

	SkipSSLValidation bool

	isPersistent bool

	originalCfHomeDir string
	currentCfHomeDir  string
}

func NewTestSetup(config config.Config) *ReproducibleTestSetup {
	node := ginkgoconfig.GinkgoConfig.ParallelNode
	timeTag := time.Now().Format("2006_01_02-15h04m05.999s")

	regUser := fmt.Sprintf("%s-USER-%d-%s", config.NamePrefix, node, timeTag)
	regUserPass := "meow"

	if config.UseExistingUser {
		regUser = config.ExistingUser
		regUserPass = config.ExistingUserPassword
	}
	if config.ConfigurableTestPassword != "" {
		regUserPass = config.ConfigurableTestPassword
	}

	return &ReproducibleTestSetup{
		config: config,

		shortTimeout: config.ScaledTimeout(1 * time.Minute),
		longTimeout:  config.ScaledTimeout(5 * time.Minute),

		quotaDefinitionName:                  generator.PrefixedRandomName(config.NamePrefix, "QUOTA"),
		quotaDefinitionTotalMemoryLimit:      "10G",
		quotaDefinitionInstanceMemoryLimit:   "-1",
		quotaDefinitionRoutesLimit:           "1000",
		quotaDefinitionAppInstanceLimit:      "-1",
		quotaDefinitionServiceInstanceLimit:  "100",
		quotaDefinitionAllowPaidServicesFlag: "--allow-paid-service-plans",

		organizationName: generator.PrefixedRandomName(config.NamePrefix, "ORG"),
		spaceName:        generator.PrefixedRandomName(config.NamePrefix, "SPACE"),

		regularUserUsername: regUser,
		regularUserPassword: regUserPass,

		adminUserUsername: config.AdminUser,
		adminUserPassword: config.AdminPassword,

		isPersistent: false,
	}
}

func NewPersistentAppTestSetup(config config.Config) *ReproducibleTestSetup {
	baseContext := NewTestSetup(config)

	baseContext.quotaDefinitionName = config.PersistentAppQuotaName
	baseContext.organizationName = config.PersistentAppOrg
	baseContext.spaceName = config.PersistentAppSpace
	baseContext.isPersistent = true

	return baseContext
}

func NewRunawayAppTestSetup(config config.Config) *ReproducibleTestSetup {
	baseContext := NewTestSetup(config)
	baseContext.quotaDefinitionTotalMemoryLimit = RUNAWAY_QUOTA_MEM_LIMIT

	return baseContext
}

func (context ReproducibleTestSetup) ShortTimeout() time.Duration {
	return context.shortTimeout
}

func (context ReproducibleTestSetup) LongTimeout() time.Duration {
	return context.longTimeout
}

func (context *ReproducibleTestSetup) Setup() {
	AsUser(context.AdminUserContext(), context.shortTimeout, func() {
		args := []string{
			"create-quota",
			context.quotaDefinitionName,
			"-m", context.quotaDefinitionTotalMemoryLimit,
			"-i", context.quotaDefinitionInstanceMemoryLimit,
			"-r", context.quotaDefinitionRoutesLimit,
			"-s", context.quotaDefinitionServiceInstanceLimit,
			"-a", context.quotaDefinitionAppInstanceLimit,
			context.quotaDefinitionAllowPaidServicesFlag,
		}

		EventuallyWithOffset(1, cf.Cf(args...), context.shortTimeout).Should(Exit(0))

		if !context.config.UseExistingUser {
			createUserCmd := cf.Cf("create-user", context.regularUserUsername, context.regularUserPassword)
			EventuallyWithOffset(1, createUserCmd, context.shortTimeout).Should(Exit())
			if createUserCmd.ExitCode() != 0 {
				ExpectWithOffset(1, createUserCmd.Out).To(Say("scim_resource_already_exists"))
			}
		}

		EventuallyWithOffset(1, cf.Cf("create-org", context.organizationName), context.shortTimeout).Should(Exit(0))
		EventuallyWithOffset(1, cf.Cf("set-quota", context.organizationName, context.quotaDefinitionName), context.shortTimeout).Should(Exit(0))

		EventuallyWithOffset(1, cf.Cf("create-space", "-o", context.organizationName, context.spaceName), context.shortTimeout).Should(Exit(0))
		EventuallyWithOffset(1, cf.Cf("set-space-role", context.regularUserUsername, context.organizationName, context.spaceName, "SpaceManager"), context.shortTimeout).Should(Exit(0))
		EventuallyWithOffset(1, cf.Cf("set-space-role", context.regularUserUsername, context.organizationName, context.spaceName, "SpaceDeveloper"), context.shortTimeout).Should(Exit(0))
		EventuallyWithOffset(1, cf.Cf("set-space-role", context.regularUserUsername, context.organizationName, context.spaceName, "SpaceAuditor"), context.shortTimeout).Should(Exit(0))
	})

	context.originalCfHomeDir, context.currentCfHomeDir = context.RegularUserContext().SetCfHomeDir()
	context.RegularUserContext().Login(context.shortTimeout)

	context.RegularUserContext().TargetSpace(context.shortTimeout)
}

func (context *ReproducibleTestSetup) Teardown() {
	context.RegularUserContext().Logout(context.shortTimeout)
	context.RegularUserContext().UnsetCfHomeDir(context.originalCfHomeDir, context.currentCfHomeDir)
	AsUser(context.AdminUserContext(), context.shortTimeout, func() {
		if !context.config.ShouldKeepUser {
			EventuallyWithOffset(1, cf.Cf("delete-user", "-f", context.regularUserUsername), context.shortTimeout).Should(Exit(0))
		}

		if !context.isPersistent {
			EventuallyWithOffset(1, cf.Cf("delete-org", "-f", context.organizationName), context.shortTimeout).Should(Exit(0))
			EventuallyWithOffset(1, cf.Cf("delete-quota", "-f", context.quotaDefinitionName), context.shortTimeout).Should(Exit(0))
		}
	})
}

func (context *ReproducibleTestSetup) AdminUserContext() UserContext {
	return NewUserContext(
		context.config.ApiEndpoint,
		context.adminUserUsername,
		context.adminUserPassword,
		"",
		"",
		context.SkipSSLValidation,
	)
}

func (context *ReproducibleTestSetup) RegularUserContext() UserContext {
	return NewUserContext(
		context.config.ApiEndpoint,
		context.regularUserUsername,
		context.regularUserPassword,
		context.organizationName,
		context.spaceName,
		context.SkipSSLValidation,
	)
}

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

const RUNAWAY_QUOTA_MEM_LIMIT = "99999G"

type ConfiguredContext struct {
	config config.Config

	shortTimeout time.Duration
	longTimeout  time.Duration

	organizationName string
	spaceName        string

	quotaDefinitionName string

	regularUserUsername string
	regularUserPassword string

	adminUserUsername string
	adminUserPassword string

	SkipSSLValidation bool

	isPersistent bool

	originalCfHomeDir string
	currentCfHomeDir  string
}

func NewContext(config config.Config) *ConfiguredContext {
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

	return &ConfiguredContext{
		config: config,

		shortTimeout: config.ScaledTimeout(1 * time.Minute),
		longTimeout:  config.ScaledTimeout(5 * time.Minute),

		quotaDefinitionName: generator.PrefixedRandomName(config.NamePrefix, "QUOTA"),

		organizationName: generator.PrefixedRandomName(config.NamePrefix, "ORG"),
		spaceName:        generator.PrefixedRandomName(config.NamePrefix, "SPACE"),

		regularUserUsername: regUser,
		regularUserPassword: regUserPass,

		adminUserUsername: config.AdminUser,
		adminUserPassword: config.AdminPassword,

		isPersistent: false,
	}
}

func NewPersistentAppContext(config config.Config) *ConfiguredContext {
	baseContext := NewContext(config)

	baseContext.quotaDefinitionName = config.PersistentAppQuotaName
	baseContext.organizationName = config.PersistentAppOrg
	baseContext.spaceName = config.PersistentAppSpace
	baseContext.isPersistent = true

	return baseContext
}

func (context ConfiguredContext) ShortTimeout() time.Duration {
	return context.shortTimeout
}

func (context ConfiguredContext) LongTimeout() time.Duration {
	return context.longTimeout
}

func (context *ConfiguredContext) Setup() {
	AsUser(context.AdminUserContext(), context.shortTimeout, func() {
		args := []string{
			"create-quota",
			context.quotaDefinitionName,
			"-m", "10G",
			"-r", "1000",
			"-s", "100",
			"--allow-paid-service-plans",
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

	context.originalCfHomeDir, context.currentCfHomeDir = InitiateUserContext(context.RegularUserContext(), context.shortTimeout)
	TargetSpace(context.RegularUserContext(), context.shortTimeout)
}

func (context *ConfiguredContext) SetRunawayQuota() {
	AsUser(context.AdminUserContext(), context.shortTimeout, func() {
		EventuallyWithOffset(1, cf.Cf("update-quota", context.quotaDefinitionName, "-m", RUNAWAY_QUOTA_MEM_LIMIT, "-i=-1"), context.shortTimeout).Should(Exit(0))
	})
}

func (context *ConfiguredContext) Teardown() {
	RestoreUserContext(context.RegularUserContext(), context.shortTimeout, context.originalCfHomeDir, context.currentCfHomeDir)
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

func (context *ConfiguredContext) AdminUserContext() UserContext {
	return NewUserContext(
		context.config.ApiEndpoint,
		context.adminUserUsername,
		context.adminUserPassword,
		"",
		"",
		context.SkipSSLValidation,
	)
}

func (context *ConfiguredContext) RegularUserContext() UserContext {
	return NewUserContext(
		context.config.ApiEndpoint,
		context.regularUserUsername,
		context.regularUserPassword,
		context.organizationName,
		context.spaceName,
		context.SkipSSLValidation,
	)
}

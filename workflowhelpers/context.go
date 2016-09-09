package workflowhelpers

import (
	"fmt"
	"time"

	ginkgoconfig "github.com/onsi/ginkgo/config"

	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
)

type ReproducibleTestSetup struct {
	config config.Config

	shortTimeout time.Duration
	longTimeout  time.Duration

	organizationName string
	spaceName        string

	testSpace *TestSpace

	regularUserContext UserContext
	adminUserContext   UserContext

	SkipSSLValidation bool

	isPersistent bool

	originalCfHomeDir string
	currentCfHomeDir  string
}

func NewTestSetup(config config.Config) *ReproducibleTestSetup {
	testSpace := NewRegularTestSpace(config)
	return newBaseTestSetup(config, testSpace)
}

func NewPersistentAppTestSetup(config config.Config) *ReproducibleTestSetup {
	testSpace := NewPersistentAppTestSpace(config)
	baseContext := newBaseTestSetup(config, testSpace)
	baseContext.isPersistent = true

	return baseContext
}

func NewRunawayAppTestSetup(config config.Config) *ReproducibleTestSetup {
	testSpace := NewRunawayAppTestSpace(config)
	return newBaseTestSetup(config, testSpace)
}

func newBaseTestSetup(config config.Config, testSpace *TestSpace) *ReproducibleTestSetup {
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

		organizationName: generator.PrefixedRandomName(config.NamePrefix, "ORG"),
		spaceName:        generator.PrefixedRandomName(config.NamePrefix, "SPACE"),

		regularUserContext: NewUserContext(config.ApiEndpoint, regUser, regUserPass, testSpace, config.SkipSSLValidation),
		adminUserContext:   NewUserContext(config.ApiEndpoint, config.AdminUser, config.AdminPassword, nil, config.SkipSSLValidation),

		isPersistent: false,
		testSpace:    testSpace,
	}
}

func (context ReproducibleTestSetup) ShortTimeout() time.Duration {
	return context.shortTimeout
}

func (context ReproducibleTestSetup) LongTimeout() time.Duration {
	return context.longTimeout
}

func (context *ReproducibleTestSetup) Setup() {
	AsUser(context.AdminUserContext(), context.shortTimeout, func() {
		context.RegularUserContext().CreateUser(context.shortTimeout)
		context.testSpace.InstantiateRemotely()
		context.RegularUserContext().AddUserToSpace(context.shortTimeout)
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
			context.RegularUserContext().DeleteUser(context.shortTimeout)
		}

		if !context.isPersistent {
			context.testSpace.Destroy()
		}
	})
}

func (context *ReproducibleTestSetup) AdminUserContext() UserContext {
	return context.adminUserContext
}

func (context *ReproducibleTestSetup) RegularUserContext() UserContext {
	return context.regularUserContext
}

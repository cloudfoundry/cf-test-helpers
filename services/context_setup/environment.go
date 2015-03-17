package context_setup

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
    "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

type TestEnvironment interface {
    Context() SuiteContext
    ShortTimeout() time.Duration
    LongTimeout() time.Duration
    AdminUserContext() cf.UserContext
    RegularUserContext() cf.UserContext
    BeforeEach()
    AfterEach()
}

type testEnvironment struct {
    context SuiteContext
    shortTimeout, longTimeout time.Duration
    originalCfHomeDir, currentCfHomeDir string
    adminUserContext, regularUserContext cf.UserContext
}

func NewTestEnvironment(context SuiteContext) testEnvironment {
    return testEnvironment{
        context: context,
        shortTimeout: context.ScaledTimeout(1 * time.Minute),
        longTimeout: context.ScaledTimeout(5 * time.Minute),
    }
}

func (e testEnvironment) Context() SuiteContext {
    return e.context
}

func (e testEnvironment) ShortTimeout() time.Duration {
    return e.shortTimeout
}

func (e testEnvironment) LongTimeout() time.Duration {
    return e.longTimeout
}

func (e testEnvironment) AdminUserContext() cf.UserContext {
    return e.adminUserContext
}

func (e testEnvironment) RegularUserContext() cf.UserContext {
    return e.regularUserContext
}

func (e testEnvironment) BeforeEach() {
    //TODO: move to constructor?
    e.adminUserContext = e.context.AdminUserContext()
    e.regularUserContext = e.context.RegularUserContext()

    e.context.Setup()

    cf.AsUser(e.adminUserContext, func() {
        e.setUpSpaceWithUserAccess(e.regularUserContext)
    })

    e.originalCfHomeDir, e.currentCfHomeDir = cf.InitiateUserContext(e.regularUserContext)
    cf.TargetSpace(e.regularUserContext)
}

func (e testEnvironment) AfterEach() {
    cf.RestoreUserContext(e.regularUserContext, e.originalCfHomeDir, e.currentCfHomeDir)

    e.context.Teardown()
}

func (e testEnvironment) setUpSpaceWithUserAccess(uc cf.UserContext) {
	runner.NewCmdRunner(cf.Cf("create-space", "-o", uc.Org, uc.Space), e.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceManager"), e.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceDeveloper"), e.shortTimeout).Run()
    runner.NewCmdRunner(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceAuditor"), e.shortTimeout).Run()
}

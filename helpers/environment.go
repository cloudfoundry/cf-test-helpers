package helpers

import (
    "time"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
)

type SuiteContext interface {
	Setup()
	Teardown()
	SetRunawayQuota()

	AdminUserContext() cf.UserContext
	RegularUserContext() cf.UserContext

    ShortTimeout() time.Duration
    LongTimeout() time.Duration
}

type Environment struct {
	context           SuiteContext
	originalCfHomeDir string
	currentCfHomeDir  string
}

func NewEnvironment(context SuiteContext) *Environment {
	return &Environment{context: context}
}

func (e *Environment) Setup() {
	e.context.Setup()

	cf.AsUser(e.context.AdminUserContext(), e.context.ShortTimeout(), func() {
		setUpSpaceWithUserAccess(e.context.RegularUserContext())
	})

	e.originalCfHomeDir, e.currentCfHomeDir = cf.InitiateUserContext(e.context.RegularUserContext(), e.context.ShortTimeout())
	cf.TargetSpace(e.context.RegularUserContext(), e.context.ShortTimeout())
}

func (e *Environment) Teardown() {
	cf.RestoreUserContext(e.context.RegularUserContext(), e.context.ShortTimeout(), e.originalCfHomeDir, e.currentCfHomeDir)

	e.context.Teardown()
}

func setUpSpaceWithUserAccess(uc cf.UserContext) {
	spaceSetupTimeout := 30.0
	Expect(cf.Cf("create-space", "-o", uc.Org, uc.Space).Wait(spaceSetupTimeout)).To(Exit(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceManager").Wait(spaceSetupTimeout)).To(Exit(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceDeveloper").Wait(spaceSetupTimeout)).To(Exit(0))
	Expect(cf.Cf("set-space-role", uc.Username, uc.Org, uc.Space, "SpaceAuditor").Wait(spaceSetupTimeout)).To(Exit(0))
}

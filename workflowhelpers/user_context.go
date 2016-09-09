package workflowhelpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/commandstarter"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	workflowhelpersinternal "github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers/internal"
	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

type UserContext struct {
	ApiUrl    string
	Username  string
	Password  string
	TestSpace *TestSpace
	Org       string // These are still used by CATS
	Space     string

	SkipSSLValidation bool
	CommandStarter    internal.Starter
}

func NewUserContext(apiUrl, username, password string, testSpace *TestSpace, skipSSLValidation bool) UserContext {
	var org, space string
	if testSpace != nil {
		org = testSpace.OrganizationName
		space = testSpace.SpaceName
	}
	return UserContext{
		ApiUrl:            apiUrl,
		Username:          username,
		Password:          password,
		TestSpace:         testSpace,
		Org:               org,
		Space:             space,
		SkipSSLValidation: skipSSLValidation,
		CommandStarter:    commandstarter.NewCommandStarter(),
	}
}

func (uc UserContext) Login(timeout time.Duration) {
	args := []string{"api", uc.ApiUrl}
	if uc.SkipSSLValidation {
		args = append(args, "--skip-ssl-validation")
	}
	session := internal.Cf(uc.CommandStarter, args...)
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))

	session = workflowhelpersinternal.CfAuth(uc.CommandStarter, uc.Username, uc.Password)
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))
}

func (uc UserContext) SetCfHomeDir() (string, string) {
	originalCfHomeDir := os.Getenv("CF_HOME")
	currentCfHomeDir, err := ioutil.TempDir("", fmt.Sprintf("cf_home_%d", ginkgoconfig.GinkgoConfig.ParallelNode))
	if err != nil {
		panic("Error: could not create temporary home directory: " + err.Error())
	}

	os.Setenv("CF_HOME", currentCfHomeDir)
	return originalCfHomeDir, currentCfHomeDir
}

func (uc UserContext) TargetSpace(timeout time.Duration) {
	if uc.TestSpace != nil {
		var session *Session
		session = internal.Cf(uc.CommandStarter, "target", "-o", uc.TestSpace.OrganizationName, "-s", uc.TestSpace.SpaceName)
		EventuallyWithOffset(1, session, timeout).Should(Exit(0))
	}
}

func (uc UserContext) AddUserToSpace(timeout time.Duration) {
	spaceManager := internal.Cf(uc.CommandStarter, "set-space-role", uc.Username, uc.TestSpace.OrganizationName, uc.TestSpace.SpaceName, "SpaceManager")
	EventuallyWithOffset(1, spaceManager, timeout).Should(Exit(0))

	spaceDeveloper := internal.Cf(uc.CommandStarter, "set-space-role", uc.Username, uc.TestSpace.OrganizationName, uc.TestSpace.SpaceName, "SpaceDeveloper")
	EventuallyWithOffset(1, spaceDeveloper, timeout).Should(Exit(0))

	spaceAuditor := internal.Cf(uc.CommandStarter, "set-space-role", uc.Username, uc.TestSpace.OrganizationName, uc.TestSpace.SpaceName, "SpaceAuditor")
	EventuallyWithOffset(1, spaceAuditor, timeout).Should(Exit(0))
}

func (uc UserContext) Logout(timeout time.Duration) {
	session := internal.Cf(uc.CommandStarter, "logout")
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))
}

func (uc UserContext) UnsetCfHomeDir(originalCfHomeDir, currentCfHomeDir string) {
	os.Setenv("CF_HOME", originalCfHomeDir)
	os.RemoveAll(currentCfHomeDir)
}

func (uc UserContext) CreateUser(timeout time.Duration) {
	session := internal.Cf(uc.CommandStarter, "create-user", uc.Username, uc.Password)
	EventuallyWithOffset(1, session, timeout).Should(Exit())
	if session.ExitCode() != 0 {
		ExpectWithOffset(1, session.Out).Should(Say("scim_resource_already_exists"))
	}
}

func (uc UserContext) DeleteUser(timeout time.Duration) {
	session := internal.Cf(uc.CommandStarter, "delete-user", "-f", uc.Username)
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))
}

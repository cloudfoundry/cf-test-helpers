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
	. "github.com/onsi/gomega/gexec"
)

type UserContext struct {
	ApiUrl   string
	Username string
	Password string
	Org      string
	Space    string

	SkipSSLValidation bool
	CommandStarter    internal.Starter
}

func NewUserContext(apiUrl, username, password, org, space string, skipSSLValidation bool) UserContext {
	return UserContext{
		ApiUrl:            apiUrl,
		Username:          username,
		Password:          password,
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
	if uc.Org != "" {
		var session *Session
		if uc.Space != "" {
			session = internal.Cf(uc.CommandStarter, "target", "-o", uc.Org, "-s", uc.Space)
		} else {
			session = internal.Cf(uc.CommandStarter, "target", "-o", uc.Org)
		}
		EventuallyWithOffset(1, session, timeout).Should(Exit(0))
	}
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
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))
}

func (uc UserContext) Destroy(timeout time.Duration) {
	session := internal.Cf(uc.CommandStarter, "delete-user", "-f", uc.Username)
	EventuallyWithOffset(1, session, timeout).Should(Exit(0))
}

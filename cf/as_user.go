package cf

import (
	"fmt"
	"io/ioutil"
	"os"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
    "time"
)

var AsUser = func(userContext UserContext, timeout time.Duration, actions func()) {
	originalCfHomeDir, currentCfHomeDir := InitiateUserContext(userContext, timeout)
	defer func() {
		RestoreUserContext(userContext, timeout, originalCfHomeDir, currentCfHomeDir)
	}()

	TargetSpace(userContext, timeout)

	actions()
}

func InitiateUserContext(userContext UserContext, timeout time.Duration) (originalCfHomeDir, currentCfHomeDir string) {
	originalCfHomeDir = os.Getenv("CF_HOME")
	currentCfHomeDir, err := ioutil.TempDir("", fmt.Sprintf("cf_home_%d", ginkgoconfig.GinkgoConfig.ParallelNode))

	if err != nil {
		panic("Error: could not create temporary home directory: " + err.Error())
	}

	os.Setenv("CF_HOME", currentCfHomeDir)

	cfSetApiArgs := []string{"api", userContext.ApiUrl}
	if userContext.SkipSSLValidation {
		cfSetApiArgs = append(cfSetApiArgs, "--skip-ssl-validation")
	}

	Expect(Cf(cfSetApiArgs...).Wait(timeout)).To(Exit(0))

	Expect(Cf("auth", userContext.Username, userContext.Password).Wait(timeout)).To(Exit(0))

	return
}

func TargetSpace(userContext UserContext, timeout time.Duration) {
	if userContext.Org != "" {
		if userContext.Space != "" {
			Expect(Cf("target", "-o", userContext.Org, "-s", userContext.Space).Wait(timeout)).To(Exit(0))
		} else {
			Expect(Cf("target", "-o", userContext.Org).Wait(timeout)).To(Exit(0))
		}
	}
}

func RestoreUserContext(_ UserContext, timeout time.Duration, originalCfHomeDir, currentCfHomeDir string) {
	Expect(Cf("logout").Wait(timeout)).To(Exit(0))
	os.Setenv("CF_HOME", originalCfHomeDir)
	os.RemoveAll(currentCfHomeDir)
}

package cf

import (
	"fmt"
	"io/ioutil"
	"os"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"
)

func AsUser(user UserContext, actions func()) {
	originalCfHomeDir := os.Getenv("CF_HOME")
	cfHomeDir, err := ioutil.TempDir("", fmt.Sprintf("cf_home_%d", ginkgoconfig.GinkgoConfig.ParallelNode))

	if err != nil {
		panic("Error: could not create temporary home directory: " + err.Error())
	}

	os.Setenv("CF_HOME", cfHomeDir)

	defer func() {
		Expect(Cf("logout")).To(ExitWith(0))
		os.Setenv("CF_HOME", originalCfHomeDir)
		os.RemoveAll(cfHomeDir)
	}()

	Expect(Cf("api", user.ApiUrl)).To(ExitWith(0))
	Expect(Cf("auth", user.Username, user.Password)).To(ExitWith(0))
	Expect(Cf("target", "-o", user.Org, "-s", user.Space)).To(ExitWith(0))

	actions()
}

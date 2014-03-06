package cf

import (
	"fmt"
	"io/ioutil"
	"os"

	ginkgoconfig "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/vito/cmdtest/matchers"
)

func AsUser(user User, actions func() error) error {
	originalCfHomeDir := os.Getenv("CF_HOME")
	cfHomeDir, err := ioutil.TempDir("", fmt.Sprintf("cf_home_%d", ginkgoconfig.GinkgoConfig.ParallelNode))

	if err != nil {
		return err
	}

	os.Setenv("CF_HOME", cfHomeDir)

	defer func() {
		Expect(Cf("logout")).To(ExitWith(0))
		os.Setenv("CF_HOME", originalCfHomeDir)
		os.RemoveAll(cfHomeDir)
	}()

	if Expect(Cf("login", user.Username, user.Password)).To(ExitWith(0)) {
		return actions()
	}

	return nil
}

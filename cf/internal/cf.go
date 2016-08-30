package cfinternal

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/commandstarter"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	"github.com/onsi/gomega/gexec"
)

func Cf(cmdStarter internal.Starter, args ...string) *gexec.Session {
	reporter := commandstarter.NewDefaultReporter()
	request, err := cmdStarter.Start(reporter, "cf", args...)
	if err != nil {
		panic(err)
	}
	return request
}

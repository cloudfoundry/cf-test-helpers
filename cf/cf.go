package cf

import (
	"os"
	"strconv"

	"github.com/onsi/ginkgo/config"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"
	"github.com/vito/cmdtest"
)

var Cf = func(args ...string) *cmdtest.Session {
	trace_file := os.Getenv("CF_TRACE_BASENAME")
	if trace_file != "" {
		os.Setenv("CF_TRACE", trace_file+strconv.Itoa(parallelNode())+".txt")
	}

	return runner.Run("gcf", args...)
}

func parallelNode() int {
	return config.GinkgoConfig.ParallelNode
}

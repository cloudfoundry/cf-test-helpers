package generator

import (
	"strconv"

	"github.com/adam-hanna/randomstrings"
	"github.com/onsi/ginkgo/config"
)

func randomName() string {
	str, err := randomstrings.GenerateRandomString(20)
	if err != nil {
		panic(err)
	}

	return str
}

func PrefixedRandomName(prefixName, resourceName string) string {
	return prefixName + "-" + strconv.Itoa(config.GinkgoConfig.ParallelNode) + "-" + resourceName + "-" + randomName()
}

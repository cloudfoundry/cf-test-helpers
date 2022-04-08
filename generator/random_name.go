package generator

import (
	"crypto/rand"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"strconv"
)

func randomName() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic("random failed:" + err.Error())
	}

	return fmt.Sprintf("%x", b)
}

func PrefixedRandomName(prefixName, resourceName string) string {
	return prefixName + "-" + strconv.Itoa(ginkgo.GinkgoParallelProcess()) + "-" + resourceName + "-" + randomName()
}

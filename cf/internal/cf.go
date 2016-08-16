package cfinternal

import (
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func Cf(cmdStarter starter, args ...string) *gexec.Session {
	request, err := cmdStarter.Start("cf", args...)
	Expect(err).NotTo(HaveOccurred())

	return request
}

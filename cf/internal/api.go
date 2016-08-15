package cfinternal

import (
	"encoding/json"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

type starter interface {
	Start(string, ...string) (*Session, error)
}

func ApiRequest(cmdStarter starter, method, endpoint string, response interface{}, timeout time.Duration, data ...string) {
	args := []string{
		"curl",
		endpoint,
		"-X", method,
	}

	dataArg := strings.Join(data, "")
	if len(dataArg) > 0 {
		args = append(args, "-d", dataArg)
	}

	request, err := cmdStarter.Start("cf", args...)
	Expect(err).NotTo(HaveOccurred())

	request.Wait(timeout)
	Expect(request).To(Exit(0))

	if response != nil {
		err := json.Unmarshal(request.Out.Contents(), response)
		Expect(err).ToNot(HaveOccurred())
	}
}

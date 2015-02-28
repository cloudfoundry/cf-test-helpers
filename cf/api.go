package cf

import (
	"encoding/json"
	"strings"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

var CF_API_TIMEOUT = 30 * time.Second

type GenericResource struct {
	Metadata struct {
		Guid string `json:"guid"`
	} `json:"metadata"`
}

type QueryResponse struct {
	Resources []GenericResource `struct:"resources"`
}

var ApiRequest = func(method, endpoint string, response interface{}, data ...string) {
	request := Cf(
		"curl",
		endpoint,
		"-X", method,
		"-d", strings.Join(data, ""),
	)

	ExecWithTimeout(request, CF_API_TIMEOUT)

	if response != nil {
		err := json.Unmarshal(request.Out.Contents(), response)
		Expect(err).ToNot(HaveOccurred())
	}
}

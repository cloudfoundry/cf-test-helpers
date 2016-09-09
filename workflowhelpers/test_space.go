package workflowhelpers

import (
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/commandstarter"
	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
)

const RUNAWAY_QUOTA_MEM_LIMIT = "99999G"

type TestSpace struct {
	QuotaDefinitionName                  string
	OrganizationName                     string
	SpaceName                            string
	IsPersistent                         bool
	QuotaDefinitionTotalMemoryLimit      string
	QuotaDefinitionInstanceMemoryLimit   string
	QuotaDefinitionRoutesLimit           string
	QuotaDefinitionAppInstanceLimit      string
	QuotaDefinitionServiceInstanceLimit  string
	QuotaDefinitionAllowPaidServicesFlag string
	CommandStarter                       internal.Starter
	Timeout                              time.Duration
}

func NewRegularTestSpace(config config.Config) *TestSpace {
	testSpace := &TestSpace{
		QuotaDefinitionName:                  generator.PrefixedRandomName(config.NamePrefix, "QUOTA"),
		QuotaDefinitionTotalMemoryLimit:      "10G",
		QuotaDefinitionInstanceMemoryLimit:   "-1",
		QuotaDefinitionRoutesLimit:           "1000",
		QuotaDefinitionAppInstanceLimit:      "-1",
		QuotaDefinitionServiceInstanceLimit:  "100",
		QuotaDefinitionAllowPaidServicesFlag: "--allow-paid-service-plans",
		OrganizationName:                     generator.PrefixedRandomName(config.NamePrefix, "ORG"),
		SpaceName:                            generator.PrefixedRandomName(config.NamePrefix, "SPACE"),
		CommandStarter:                       commandstarter.NewCommandStarter(),
		Timeout:                              config.ScaledTimeout(1 * time.Minute),
	}
	return testSpace
}

func NewPersistentAppTestSpace(config config.Config) *TestSpace {
	baseTestSpace := NewRegularTestSpace(config)
	baseTestSpace.IsPersistent = true
	baseTestSpace.OrganizationName = config.PersistentAppOrg
	baseTestSpace.SpaceName = config.PersistentAppSpace
	baseTestSpace.QuotaDefinitionName = config.PersistentAppQuotaName
	return baseTestSpace
}

func NewRunawayAppTestSpace(config config.Config) *TestSpace {
	baseTestSpace := NewRegularTestSpace(config)
	baseTestSpace.QuotaDefinitionTotalMemoryLimit = RUNAWAY_QUOTA_MEM_LIMIT
	return baseTestSpace
}

func (ts *TestSpace) InstantiateRemotely() {
	args := []string{
		"create-quota",
		ts.QuotaDefinitionName,
		"-m", ts.QuotaDefinitionTotalMemoryLimit,
		"-i", ts.QuotaDefinitionInstanceMemoryLimit,
		"-r", ts.QuotaDefinitionRoutesLimit,
		"-a", ts.QuotaDefinitionAppInstanceLimit,
		"-s", ts.QuotaDefinitionServiceInstanceLimit,
		ts.QuotaDefinitionAllowPaidServicesFlag,
	}

	createQuota := internal.Cf(ts.CommandStarter, args...)
	EventuallyWithOffset(1, createQuota, ts.Timeout).Should(Exit(0))

	createOrg := internal.Cf(ts.CommandStarter, "create-org", ts.OrganizationName)
	EventuallyWithOffset(1, createOrg, ts.Timeout).Should(Exit(0))

	setQuota := internal.Cf(ts.CommandStarter, "set-quota", ts.OrganizationName, ts.QuotaDefinitionName)
	EventuallyWithOffset(1, setQuota, ts.Timeout).Should(Exit(0))

	createSpace := internal.Cf(ts.CommandStarter, "create-space", "-o", ts.OrganizationName, ts.SpaceName)
	EventuallyWithOffset(1, createSpace, ts.Timeout).Should(Exit(0))
}

func (ts *TestSpace) Destroy() {
	if !ts.IsPersistent {
		deleteOrg := internal.Cf(ts.CommandStarter, "delete-org", "-f", ts.OrganizationName)
		EventuallyWithOffset(1, deleteOrg, ts.Timeout).Should(Exit(0))

		deleteQuota := internal.Cf(ts.CommandStarter, "delete-quota", "-f", ts.QuotaDefinitionName)
		EventuallyWithOffset(1, deleteQuota, ts.Timeout).Should(Exit(0))
	}
}

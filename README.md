cf-test-helpers
===============

Go utilities for running tests against Cloud Foundry

Included Tools:
- Execute [CF CLI](https://github.com/cloudfoundry/cli) commands
  - Isolated user contexts for wrapping CF commands
  - Curl CF endpoints
- Random user name generator
- Thin wrapper around curl (in `cf-test-helpers/runner`)

## Observing cf commands (optional)

You can observe every `cf` command the helpers run — its redacted identity,
elapsed wall-clock time, and exit code — by registering an observer once, e.g.
at suite start. This is opt-in: with no observer registered, behavior is
unchanged.

```go
import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/cf"
	. "github.com/onsi/ginkgo/v2"
)

// commandCounter tracks call counts and slowest elapsed time per command.
// Safe for concurrent use — CommandCompleted fires from a background goroutine.

type commandCounter struct {
	mu     sync.Mutex
	counts map[string]int
	slow   map[string]time.Duration
}

func (c *commandCounter) CommandStarted(string, time.Time) {}
func (c *commandCounter) CommandCompleted(cmd string, elapsed time.Duration, _ int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[cmd]++
	if elapsed > c.slow[cmd] {
		c.slow[cmd] = elapsed
	}
}

// In your suite file:
var counter = &commandCounter{
	counts: map[string]int{},
	slow:   map[string]time.Duration{},
}

var _ = BeforeSuite(func() { cf.RegisterObserver(counter) })
var _ = AfterSuite(func() {
	cf.UnregisterObserver()
	for cmd, n := range counter.counts {
		fmt.Printf("%s: %d calls, slowest %s\n", cmd, n, counter.slow[cmd])
	}
})
```

Notes:
- `redactedArgs` is the command's arguments joined into one string, with any
  library-known secret (e.g. the `CfRedact` secret) already redacted. It never
  contains a library-known secret in the clear.
- `CommandStarted` fires for every command. `CommandCompleted` fires from the
  library's own goroutine when the process terminates — independent of whether
  the caller waits on the session. A command whose process never terminates
  (or that outlives the test program) produces a start event and no completion
  event.
- Registration is process-global and concurrency-safe; set it once.
- Registering an observer does not suppress the usual human-readable command
  line the helpers print; you get both the console line and the observer
  callbacks (redaction applies to both).
- Your live `*gexec.Session` and existing `Eventually(session).Should(Exit(0))`
  assertions are unaffected.

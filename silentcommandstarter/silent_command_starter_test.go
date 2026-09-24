package silentcommandstarter_test

import (
	"os/exec"
	"sync"
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/silentcommandstarter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type completionReporter struct {
	mu        sync.Mutex
	completed bool
	args      []string
	elapsed   time.Duration
	exitCode  int
}

func (f *completionReporter) Report(t time.Time, cmd *exec.Cmd) {}

func (f *completionReporter) ReportCompletion(cmd *exec.Cmd, elapsed time.Duration, exitCode int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.completed = true
	f.args = cmd.Args
	f.elapsed = elapsed
	f.exitCode = exitCode
}

func (f *completionReporter) snapshot() (bool, []string, time.Duration, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.completed, f.args, f.elapsed, f.exitCode
}

var _ = Describe("SilentCommandStarter", func() {
	var cmdStarter *silentcommandstarter.CommandStarter

	BeforeEach(func() {
		cmdStarter = silentcommandstarter.NewCommandStarter()
	})

	When("the reporter implements CompletionReporter", func() {
		It("reports completion with exit code and non-negative elapsed, without the caller calling Wait()", func() {
			cr := &completionReporter{}
			session, err := cmdStarter.Start(cr, "bash", "-c", "exit 5")
			Expect(err).To(Succeed())

			Eventually(func() bool {
				done, _, _, _ := cr.snapshot()
				return done
			}, 5*time.Second).Should(BeTrue())

			_, args, elapsed, exit := cr.snapshot()
			Expect(args).To(Equal([]string{"bash", "-c", "exit 5"}))
			Expect(exit).To(Equal(5))
			Expect(elapsed).To(BeNumerically(">=", 0))

			_ = session
		})
	})
})

package formatter_test

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/formatter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("CliErrorMessage", func() {
	It("returns a nice error message that includes the command that was run", func() {
		args := []string{"cf", "push"}
		command := &exec.Cmd{
			Args: args,
		}
		session := &gexec.Session{
			Command: command,
		}

		argsAsString := strings.Join(args, " ")
		helpfulMessage := fmt.Sprintf("\n>>> [ %s ] exited with an error \n", argsAsString)
		Expect(formatter.CliErrorMessage(session)).To(Equal(helpfulMessage))
	})

	It("redacts args for the `cf auth` command", func() {
		args := []string{"cf", "auth", "super", "secret"}
		command := &exec.Cmd{
			Args: args,
		}
		session := &gexec.Session{
			Command: command,
		}

		argsAsString := strings.Join(args[:2], " ")
		helpfulMessage := fmt.Sprintf("\n>>> [ %s ] exited with an error \n", argsAsString)
		Expect(formatter.CliErrorMessage(session)).To(Equal(helpfulMessage))
	})

	It("panics if the session is empty", func() {
		session := &gexec.Session{}

		defer func() {
			if r := recover(); r != nil {
				Expect(r).To(Equal("session was nil!"))
			}
		}()

		formatter.CliErrorMessage(session)
	})

	It("panics if the command is empty", func() {
		command := &exec.Cmd{}
		session := &gexec.Session{
			Command: command,
		}

		defer func() {
			if r := recover(); r != nil {
				Expect(r).To(Equal("command was nil!"))
			}
		}()

		formatter.CliErrorMessage(session)
	})
})

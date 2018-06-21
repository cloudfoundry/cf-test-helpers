package formatter

import (
	"fmt"
	"strings"

	"github.com/onsi/gomega/gexec"
)

func CliErrorMessage(session *gexec.Session) string {
	if *session == (gexec.Session{}) {
		panic("session was nil!")
	}

	if len(session.Command.Args) == 0 {
		panic("command was nil!")
	}

	var command string

	if strings.EqualFold(session.Command.Args[1], "auth") {
		command = strings.Join(session.Command.Args[:2], " ")
	} else {
		command = strings.Join(session.Command.Args, " ")
	}

	return fmt.Sprintf("\n>>> [ %s ] exited with an error \n", command)
}

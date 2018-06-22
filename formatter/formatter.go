package formatter

// We might not end up using this code
// We found a way to put the actual error message from the CLI into our Gomega annotations
// We haven't yet removed it, pending PM review
// But if you're looking at this later and don't know why it's here, just delete it

import (
	"fmt"
	"strings"
)

func CliErrorMessage(args []string) string {
	var command string

	if strings.EqualFold(args[1], "auth") {
		command = strings.Join(args[:2], " ")
	} else {
		command = strings.Join(args, " ")
	}

	return fmt.Sprintf("\n>>> [ %s ] exited with an error \n", command)
}

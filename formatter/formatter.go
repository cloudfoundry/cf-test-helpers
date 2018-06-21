package formatter

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

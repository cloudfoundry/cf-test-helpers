package helpers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
)

func AsInterceptCommand(stdinString string, env map[string]string, actions func()) {
	originalCommandInterceptor := runner.CommandInterceptor
	runner.CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
		if stdinString != "" {
			setStdin(cmd, stdinString)
		}
		if env != nil {
			setEnv(cmd, env)
		}
		return cmd
	}
	defer resetCommandInterceptor(originalCommandInterceptor)
	actions()
}

func resetCommandInterceptor(originalCommandInterceptor func(cmd *exec.Cmd) *exec.Cmd) {
	runner.CommandInterceptor = originalCommandInterceptor
}

func setStdin(cmd *exec.Cmd, stdinString string) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	defer stdin.Close()
	io.Copy(stdin, bytes.NewBufferString(stdinString))
}

func setEnv(cmd *exec.Cmd, newEnv map[string]string) {
	orgEnv := os.Environ()
	mergedEnvMap := make(map[string]string)

	for _, v := range orgEnv {
		variables := strings.SplitN(v, "=", 2)
		orgKey := variables[0]
		orgValue := variables[1]
		mergedEnvMap[orgKey] = orgValue
	}

	for k, v := range newEnv {
		mergedEnvMap[k] = v
	}

	for k, v := range mergedEnvMap {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

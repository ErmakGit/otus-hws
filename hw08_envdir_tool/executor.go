package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return -1
	}

	execCommand := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	execCommand.Stdin = os.Stdin
	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr

	for key, value := range env {
		_, exists := os.LookupEnv(key)
		if exists {
			if err := os.Unsetenv(key); err != nil {
				log.Fatalf("error unset env var %s\n", key)
			}
		}

		if value.NeedRemove {
			continue
		}

		if err := os.Setenv(key, value.Value); err != nil {
			log.Fatalf("error set env var %s\n", key)
		}
	}

	if err := execCommand.Run(); err != nil {
		log.Fatal(err)
	}

	return execCommand.ProcessState.ExitCode()
}

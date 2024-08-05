package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("unset env", func(t *testing.T) {
		err := os.Setenv("SOME_ENV", "old_value")
		require.NoError(t, err)

		env := make(Environment)
		env["SOME_ENV"] = EnvValue{
			"new_value",
			false,
		}

		RunCmd([]string{"ls"}, env)
		someEnvVal, ok := os.LookupEnv("SOME_ENV")
		if !ok {
			fmt.Println("env not found")
		}
		require.Equal(t, "new_value", someEnvVal)
	})

	t.Run("empty environment", func(t *testing.T) {
		r := RunCmd([]string{"ls"}, Environment{})
		require.Equal(t, 0, r)
	})

	t.Run("empty cmd and environment", func(t *testing.T) {
		var s []string
		r := RunCmd(s, Environment{})
		require.Equal(t, -1, r)
	})
}

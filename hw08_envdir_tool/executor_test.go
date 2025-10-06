package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Set up the command and environment set
	cmd := []string{"printenv", "TEST_VAR__"}
	env := Environment{
		"TEST_VAR__": EnvValue{
			Value:      "test",
			NeedRemove: false,
		},
	}
	res, err := RunCmd(cmd, env)
	if err != nil {
		t.Errorf("run RunCmd with full cmd failed: %v", err)
	}
	require.Equal(t, "test\n", res)

	// empty cmd
	empyCmd := make([]string, 0)
	_, err = RunCmd(empyCmd, env)
	require.Error(t, ErrNoCommand, err)

	// singl cmd
	snglCmd := []string{"pwd"}
	res, err = RunCmd(snglCmd, env)
	if err != nil {
		t.Errorf("run RunCmd with snglCmd failed: %v", err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf(" call os.Getwd() failed: %v", err)
	}
	require.Equal(t, pwd+"\n", res)
}

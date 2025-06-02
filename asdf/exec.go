package asdf

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"al.essio.dev/pkg/shellescape"
)

// ExecEnv contains everything needed to run asdf commands in a specific environment
// that is installed and pre-configured.
type ExecEnv struct {
	// Env vars that confiure asdf and are required for its operation.
	EnvVars   map[string]string

	// ShellInit is a shell command that initializes asdf in the shell session.
	// This is required because classic asdf is written in bash and we can't assume that
	// its init command is sourced in .bashrc or similar (and we don't want to modify
	// anything system-wide).
	ShellInit string
}

func (e *ExecEnv) runAsdf(args ...string) (string, error) {
	asdfCmd := []string{}
	if e.ShellInit != "" {
		asdfCmd = append(asdfCmd, e.ShellInit+" && ")
	}
	asdfCmd = append(asdfCmd, "asdf")
	escapedAsdfArgs := shellescape.QuoteCommand(args)
	asdfCmd = append(asdfCmd, escapedAsdfArgs)

	// We need to spawn a sub-shell because classic asdf is implemented in bash and
	// relies on shell features.
	cmdArgs := []string{"-c", strings.Join(asdfCmd, " ")}
	command := exec.Command("bash", cmdArgs...)
	command.Env = os.Environ()
	for k, v := range e.EnvVars {
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, v))
	}

	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w\n\nOutput:\n%s", "asdf", escapedAsdfArgs, err, output)
	}

	return string(output), nil
}

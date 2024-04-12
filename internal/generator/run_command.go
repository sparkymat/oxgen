package generator

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *Service) runCommand(workspaceFolder string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = workspaceFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running %s %s: %w", command, strings.Join(args, " "), err)
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

func (s *Service) Generate(ctx context.Context) error {
	slog.Info("running pre commands", "commands", s.Config.PreCommands)

	for _, command := range s.Config.PreCommands {
		cmd := exec.CommandContext(ctx, "sh", "-c", command)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run command '%s': %w", command, err)
		}
	}

	err := s.GenerateProject(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate project folder: %w", err)
	}

	for _, resource := range s.Config.Resources {
		err := s.GenerateResource(ctx, resource)
		if err != nil {
			return fmt.Errorf("failed to generate resource folder: %w", err)
		}
	}

	slog.Info("running post commands", "commands", s.Config.PostCommands)

	for _, command := range s.Config.PostCommands {
		cmd := exec.CommandContext(ctx, "sh", "-c", command)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run command '%s': %w", command, err)
		}
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

	slog.Info("walking through templates folder", "path", path.Join(s.Config.TemplatesFolder, "project"))

	err := filepath.WalkDir(
		path.Join(s.Config.TemplatesFolder, "project"),
		processTemplateFile(ctx, s),
	)
	if err != nil {
		return fmt.Errorf("failed to process setup templates: %w", err)
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

func processTemplateFile(ctx context.Context, s *Service) func(string, fs.DirEntry, error) error {
	return func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(filePath) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		localFolderPath := fmt.Sprintf("%s%c", path.Join(s.Config.TemplatesFolder, "project"), os.PathSeparator)
		destinationPathTemplate, _ := strings.CutPrefix(filePath, localFolderPath)
		destinationPathTemplate = strings.TrimSuffix(destinationPathTemplate, ".ot")

		destinationPath := s.ReplacePlaceholders(ctx, destinationPathTemplate)

		extention := path.Ext(destinationPath)

		destinationContents := string(content)

		if extention == ".ot" {
			destinationContents = s.ReplacePlaceholders(ctx, destinationContents)
		}

		if err := os.WriteFile(destinationPath, []byte(destinationContents), 0o644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	}
}

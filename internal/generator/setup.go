package generator

import (
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s *Service) Setup(ctx context.Context, projectName string, _ bool) error {
	err := filepath.WalkDir(path.Join("templates", "setup"), processSetupTemplateFile(ctx, s, projectName))
	if err != nil {
		return err
	}

	return nil
}

func processSetupTemplateFile(ctx context.Context, s *Service, projectName string) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		destinationPathTemplate, _ := strings.CutPrefix(path, "templates/setup/")
		destinationPathTemplate = strings.TrimSuffix(destinationPathTemplate, ".ot")

		destinationPath := s.ReplacePlaceholders(ctx, destinationPathTemplate, projectName)

		destinationContents := s.ReplacePlaceholders(ctx, string(content), projectName)

		if err := os.WriteFile(destinationPath, []byte(destinationContents), 0644); err != nil {
			return err
		}

		return nil
	}
}

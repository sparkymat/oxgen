package generator

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s *Service) Setup(ctx context.Context, projectName string, templatesFolder string, _ bool) error {
	err := filepath.WalkDir(path.Join(templatesFolder, "setup"), processSetupTemplateFile(ctx, s, projectName, templatesFolder))
	if err != nil {
		return err
	}

	return nil
}

func processSetupTemplateFile(ctx context.Context, s *Service, projectName string, templatesFolder string) func(string, fs.DirEntry, error) error {
	return func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		localFolderPath := fmt.Sprintf("%s%c", path.Join(templatesFolder, "setup"), os.PathSeparator)
		destinationPathTemplate, _ := strings.CutPrefix(filePath, localFolderPath)
		destinationPathTemplate = strings.TrimSuffix(destinationPathTemplate, ".ot")

		destinationPath := s.ReplacePlaceholders(ctx, destinationPathTemplate, projectName)

		destinationContents := s.ReplacePlaceholders(ctx, string(content), projectName)

		if err := os.WriteFile(destinationPath, []byte(destinationContents), 0644); err != nil {
			return err
		}

		return nil
	}
}

package generator

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (s *Service) GenerateProject(ctx context.Context) error {
	slog.Info("walking through project folder", "path", path.Join(s.Config.TemplatesFolder, "project"))

	projectLookupTable := generateLookupTableForProject(s.Config)

	err := filepath.WalkDir(
		path.Join(s.Config.TemplatesFolder, "project"),
		processTemplateFile(ctx, s, projectLookupTable),
	)
	if err != nil {
		return fmt.Errorf("failed to process setup templates: %w", err)
	}

	slog.Info("running post commands", "commands", s.Config.PostCommands)

	return nil
}

func processTemplateFile(ctx context.Context, s *Service, lookupTable map[string]string) func(string, fs.DirEntry, error) error {
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

		destinationPath := replacePlaceholders(ctx, lookupTable, destinationPathTemplate)
		destinationFolder := path.Dir(destinationPath)

		if err := os.MkdirAll(destinationFolder, 0o755); err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}

		destinationContents := string(content)

		ext := path.Ext(filePath)
		if ext == ".ot" {
			destinationContents = replacePlaceholders(ctx, lookupTable, destinationContents)
		}

		if err := os.WriteFile(destinationPath, []byte(destinationContents), 0o644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	}
}

package generator

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"path/filepath"
)

func (s *Service) GenerateResource(ctx context.Context, resource Resource) error {
	slog.Info("walking through resource folder", "path", path.Join(s.Config.TemplatesFolder, "resource"))

	resourceLookupTable := generateLookupTableForResource(s.Config, resource.Name)

	err := filepath.WalkDir(
		path.Join(s.Config.TemplatesFolder, "resource"),
		processTemplateFile(ctx, s, resourceLookupTable),
	)
	if err != nil {
		return fmt.Errorf("failed to process setup templates: %w", err)
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const presenterTemplate = `
package presenter

type {{ .Resource.CamelcaseSingular }} struct {
  ID string ` + "`json:\"id\"`" + `
{{range .Fields }}{{ .PresenterGoFragment }}
{{end}}  CreatedAt string ` + "`json:\"createdAt\"`" + `
  UpdatedAt string ` + "`json:\"updatedAt\"`" + `
}

func {{ .Resource.CamelcaseSingular }}FromModel(m dbx.{{ .Resource.CamelcaseSingular }}) {{ .Resource.CamelcaseSingular }} {
  item := {{ .Resource.CamelcaseSingular }}{
    ID: m.ID.String(),
    CreatedAt: m.CreatedAt.Time.Format(time.RFC3339),
    UpdatedAt: m.UpdatedAt.Time.Format(time.RFC3339),
  }

  {{range .Fields }}{{ .PresenterAssignment }}{{ end }}
  return item
}
`

func (s *Service) generatePresenter(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "handler", "api", "presenter")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure presenter folder exists: %w", err)
	}

	filename := input.Resource.UnderscoreSingular() + ".go"

	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"presenter",
		presenterTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to generate presenter method: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

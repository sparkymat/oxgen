package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

type ServiceTemplateInput struct {
	Service TemplateName
}

const serviceTemplate = `package {{ .Service.Downcase }}

func New(storageFolder string, dbx service.DatabaseProvider) *Service {
  return &Service{
    storageFolder: storageFolder,
    dbx: dbx,
  }
}

type Service struct {
  storageFolder string
  dbx service.DatabaseProvider
}
`

func (s *Service) ensureServiceExists(_ context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service.String())
	filePath := filepath.Join(folderPath, "service.go")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure service folder exists: %w", err)
	}

	if err := s.ensureFileExists(filePath, input.Service.String(), serviceTemplate, input); err != nil {
		return fmt.Errorf("failed to ensure service file exists: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", "service.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

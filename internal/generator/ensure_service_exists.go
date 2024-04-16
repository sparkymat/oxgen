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

func New(dbx service.DatabaseProvider) *Service {
  return &Service{dbx: dbx}
}

type Service struct {
  dbx service.DatabaseProvider
}
`

func (s *Service) ensureServiceExists(_ context.Context, input GenerateInput) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service)
	filePath := filepath.Join(folderPath, "service.go")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure service folder exists: %w", err)
	}

	templateInput := ServiceTemplateInput{
		Service: TemplateName(input.Service),
	}

	if err := s.ensureFileExists(filePath, input.Service, serviceTemplate, templateInput); err != nil {
		return fmt.Errorf("failed to ensure service file exists: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", "service.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

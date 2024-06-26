package generator

import (
	"context"
	"errors"
	"fmt"
	"os"
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

const mainServiceInitTemplate = `{{ .Service.String }}Service := {{ .Service.String }}.New(cfg.StorageFolder(), db)
services.{{ .Service.Capitalize }} = {{ .Service.String }}Service
`

func (s *Service) ensureServiceExists(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service.String())
	filePath := filepath.Join(folderPath, "service.go")

	serviceFileExists := true
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		serviceFileExists = false
	}

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure service folder exists: %w", err)
	}

	if err := s.ensureFileExists(filePath, input.Service.String(), serviceTemplate, input); err != nil {
		return fmt.Errorf("failed to ensure service file exists: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", "service.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	if !serviceFileExists {
		// add new Service to services interface
		if err := s.addServiceToServices(ctx, input); err != nil {
			return fmt.Errorf("failed adding service methods to interface: %w", err)
		}

		// add the new service to main.go
		if err := s.injectTemplateAboveLine(
			filepath.Join(input.WorkspaceFolder, "main.go"),
			"// End of main code generated by oxgen. DO NOT EDIT.",
			"main-service",
			mainServiceInitTemplate,
			input,
		); err != nil {
			return fmt.Errorf("failed to inject service init into main: %w", err)
		}

		if err := s.runCommand(input.WorkspaceFolder, "goimports", "-w", "main.go"); err != nil {
			return fmt.Errorf("failed running goimports: %w", err)
		}
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const servicesStructTemplate = `package internal

type Services struct {}
`

const serviceServicesStructTemplate = `
{{ .Service.Capitalize }} {{ .Service.Capitalize }}Service
`

func (s *Service) addServiceToServices(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal")
	structFilePath := filepath.Join(folderPath, "services.go")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure internal folder exists: %w", err)
	}

	if err := s.ensureFileExists(structFilePath, "services", servicesStructTemplate, input); err != nil {
		return fmt.Errorf("failed to ensure service struct file exists: %w", err)
	}

	//nolint:gomnd
	if err := s.appendTemplateToFile(ctx, structFilePath, 2, "}", "serviceEntry"+input.Service.Capitalize(), serviceServicesStructTemplate, input); err != nil {
		return fmt.Errorf("failed to append service entry to services struct: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", "services.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

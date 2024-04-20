package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const serviceIfaceTemplate = `package handler

type {{ .Service.Capitalize }}Service interface {}
`

const serviceMethodsTemplate = `
  Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params {{ .Service.String }}.Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error){{if .HasSearch }}
  Search{{ .Resource.CamelcasePlural }}(ctx context.Context, query string, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error){{end}}
  FetchRecent{{ .Resource.CamelcasePlural }}(ctx context.Context, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error) 
  Fetch{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID)(dbx.{{ .Resource.CamelcaseSingular }}, error) 
  Destroy{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID) error 
`

const updateServiceMethodsTemplate = `
  Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, {{ .UpdateGoFunctionSignatureParam }}) (dbx.{{ .Resource.CamelcaseSingular }}, error)
`

func (s *Service) addServiceMethodsToIface(
	ctx context.Context,
	input Input,
) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "handler")
	ifaceFilePath := filepath.Join(folderPath, input.Service.String()+"_service_iface.go")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure handler folder exists: %w", err)
	}

	if err := s.ensureFileExists(ifaceFilePath, input.Service.String(), serviceIfaceTemplate, input); err != nil {
		return fmt.Errorf("failed to ensure service iface file exists: %w", err)
	}

	if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "serviceMethods", serviceMethodsTemplate, input); err != nil {
		return err
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "updateServiceMethod", updateServiceMethodsTemplate, field); err != nil {
				return fmt.Errorf("failed to add update %s service method to iface file: %w", field.Name.String(), err)
			}
		}
	}

	if err := s.runCommand(folderPath, "goimports", "-w", input.Service.String()+"_service_iface.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

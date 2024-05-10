//nolint:lll,revive
package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const serviceIfaceTemplate = `package internal

type {{ .Service.Capitalize }}Service interface {}
`

const serviceMethodsIfaceTemplate = `
  Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params {{ .Service.String }}.Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error){{if .HasSearch }}
  Search{{ .Resource.CamelcasePlural }}(ctx context.Context,{{if ne .Parent nil}}parentID uuid.UUID,{{end}} query string, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error){{end}}
  FetchRecent{{ .Resource.CamelcasePlural }}(ctx context.Context,{{if ne .Parent nil}}parentID uuid.UUID,{{end}} pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error)
  Fetch{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID)(dbx.{{ .Resource.CamelcaseSingular }}, error) 
  Destroy{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID) error 
`

const updateServiceMethodsIfaceTemplate = `
  Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, {{ .UpdateGoFunctionSignatureParam }}) (dbx.{{ .Resource.CamelcaseSingular }}, error)
`

const uploadAttachmentServiceMethodsIfaceTemplate = `
  Upload{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.{{ .Resource.CamelcaseSingular }}, error)
`

func (s *Service) addServiceMethodsToIface(
	ctx context.Context,
	input Input,
) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal")
	ifaceFilePath := filepath.Join(folderPath, input.Service.String()+"_service_iface.go")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure handler folder exists: %w", err)
	}

	if err := s.ensureFileExists(ifaceFilePath, input.Service.String(), serviceIfaceTemplate, input); err != nil {
		return fmt.Errorf("failed to ensure service iface file exists: %w", err)
	}

	//nolint:gomnd
	if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "serviceMethods", serviceMethodsIfaceTemplate, input); err != nil {
		return err
	}

	for _, field := range input.Fields {
		//nolint:nestif
		if field.Updateable {
			if field.Type == FieldTypeAttachment {
				//nolint:gomnd
				if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "uploadAttachmentServiceMethod", uploadAttachmentServiceMethodsIfaceTemplate, field); err != nil {
					return fmt.Errorf("failed to add upload attachment %s service method to iface file: %w", field.Name.String(), err)
				}
			} else {
				//nolint:gomnd
				if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "updateServiceMethod", updateServiceMethodsIfaceTemplate, field); err != nil {
					return fmt.Errorf("failed to add update %s service method to iface file: %w", field.Name.String(), err)
				}
			}
		}
	}

	if err := s.runCommand(folderPath, "goimports", "-w", input.Service.String()+"_service_iface.go"); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

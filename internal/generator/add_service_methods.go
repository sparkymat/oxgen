package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const createServiceMethodTemplate = `
package {{ .Service }}

type Create{{ .Resource.CamelcaseSingular }}Params struct {
{{range .Fields }}{{if .Initial}}{{ .Field.CamelcaseSingular }} {{ .Type.GoType }}{{end}}
{{end}}
}

func (s *Service) Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error) {
  input := dbx.Create{{ .Resource.CamelcaseSingular }}Params{
{{range .Fields }}{{if .Initial}}{{ .Field.CamelcaseSingular }}: params.{{ .Field.CamelcaseSingular }},{{end}}
{{end}}
  }

  val, err := s.dbx.Create{{ .Resource.CamelcaseSingular }}(ctx, input)
  if err != nil {
    return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to create {{ .Resource.CamelcaseSingular }}: %w", err)
  }

  return val, nil
}
`

const updateServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, value {{ .Type.GoType }}) (dbx.{{ .Resource.CamelcaseSingular }}, error) {
  input := dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamelcaseSingular }}Params{
    ID: id,
    {{ .Field.CamelcaseSingular }}: value,
  }

  val, err := s.dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamelcaseSingular }}(ctx, input)
  if err != nil {
    return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to update {{ .Resource.CamelcaseSingular }} {{ .Field.CamelcaseSingular }}: %w", err)
  }

  return val, nil
}
`

func (s *Service) addServiceMethods(ctx context.Context, input GenerateInput) error {
	templateInput := TemplateInputFromGenerateInput(input)

	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service)

	// Create
	filename := fmt.Sprintf("create_%s.go", templateInput.Resource.UnderscoreSingular())
	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"create",
		createServiceMethodTemplate,
		templateInput,
	); err != nil {
		return fmt.Errorf("failed to append create service method: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	// Updates
	for _, field := range templateInput.Fields {
		if field.Updateable {
			filename := fmt.Sprintf("update_%s_%s.go", field.Resource.UnderscoreSingular(), field.Field.UnderscoreSingular())
			filePath := filepath.Join(folderPath, filename)
			if err := s.appendTemplateToFile(
				ctx,
				filePath,
				0,
				"",
				fmt.Sprintf("update_%s", field.Field.String()),
				updateServiceMethodTemplate,
				field,
			); err != nil {
				return fmt.Errorf("failed to append update service method: %w", err)
			}

			if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
				return fmt.Errorf("failed running goimports: %w", err)
			}
		}
	}

	return nil
}

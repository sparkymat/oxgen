package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const createServiceMethodTemplate = `
package {{ .Service }}

type Create{{ .Resource.CamelcaseSingular }}Params struct {
{{range .Fields }}{{if .Initial}}{{ .CreateParamsGoFragment }}{{end}}
{{end}}
}

func (s *Service) Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error) {
  input := dbx.Create{{ .Resource.CamelcaseSingular }}Params{
{{range .Fields }}{{if .Initial}}{{ .CrateAssignParamsGoFragment }},{{end}}
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

func (s *Service) Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, {{ .UpdateGoFunctionSignatureParam }}) (dbx.{{ .Resource.CamelcaseSingular }}, error) {
  {{if .NotNull}}{{else}}value := {{ .PgZeroValue }}
    if valuePtr != nil {
      value = {{ .PgValue }}
    }

    {{end}}input := dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}Params{
    ID: id,
    {{ .Name.CamelcaseSingular }}: value,
  }

  val, err := s.dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx, input)
  if err != nil {
    return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to update {{ .Resource.CamelcaseSingular }} {{ .Name.CamelcaseSingular }}: %w", err)
  }

  return val, nil
}
`

func (s *Service) addServiceMethods(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service.String())

	// Create
	filename := fmt.Sprintf("create_%s.go", input.Resource.UnderscoreSingular())
	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"create",
		createServiceMethodTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to append create service method: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	// Updates
	for _, field := range input.Fields {
		if field.Updateable {
			filename := fmt.Sprintf("update_%s_%s.go", field.Resource.UnderscoreSingular(), field.Name.UnderscoreSingular())
			filePath := filepath.Join(folderPath, filename)
			if err := s.appendTemplateToFile(
				ctx,
				filePath,
				0,
				"",
				fmt.Sprintf("update_%s", field.Name.String()),
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

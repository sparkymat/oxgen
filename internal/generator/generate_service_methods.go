//nolint:lll,revive
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
{{range .Fields }}{{if .Initial}}{{ .CreateAssignParamsGoFragment }},{{end}}
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

const uploadAttachmentServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Upload{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx context.Context, id uuid.UUID, filename string, attachmentFile io.Reader) (dbx.{{ .Resource.CamelcaseSingular }}, error) {
	folderPath := path.Join(s.storageFolder, "{{ .Resource.UnderscoreSingular }}", id.String())

	if err := os.MkdirAll(folderPath, 0o755); err != nil { //nolint:gomnd
		return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to create {{ .Resource.UnderscoreSingular }} folder. err: %w", err)
	}

	filePath := path.Join(folderPath, filename)

	data, err := io.ReadAll(attachmentFile)
	if err != nil {
		return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to read attachment file: %w", err)
	}

	//nolint:gomnd
	if err = os.WriteFile(filePath, data, 0o600); err != nil {
		return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to write attachment file: %w", err)
	}

  input := dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}Params{
    ID: id,
    {{ .Name.CamelcaseSingular }}: pgtype.Text{String: fmt.Sprintf("/{{ .Resource.UnderscoreSingular }}/%s/%s", id.String(), filename), Valid: true},
  }

	item, err := s.dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(ctx, input)
  if err != nil {
		return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to update artist photo: %w", err)
	}

	return item, nil
}
`

const searchServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Search{{ .Resource.CamelcasePlural }}(ctx context.Context,{{if ne .Parent nil}}parentID uuid.UUID,{{end}} query string, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error) {
	offset := (pageNumber - 1) * pageSize

	items, err := s.dbx.Search{{ .Resource.CamelcasePlural }}(ctx, dbx.Search{{ .Resource.CamelcasePlural }}Params{
		Query:      query,{{if ne .Parent nil}}
    ParentID: parentID,
{{end}}		PageOffset: offset,
		PageLimit:  pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search {{ .Resource.CamelcasePlural }}: %w", err)
	}

{{if eq .Parent nil}}	totalCount, err := s.dbx.CountSearched{{ .Resource.CamelcasePlural }}(ctx, query)
{{else}}  totalCount, err := s.dbx.CountSearched{{ .Resource.CamelcasePlural }}(ctx, dbx.CountSearched{{ .Resource.CamelcasePlural }}Params{
    Query: query,
    ParentID: parentID,
  })
{{end}}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch {{ .Resource.CamelcasePlural }} search count: %w", err)
	}

	return items, totalCount, nil
}
`

const recentServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) FetchRecent{{ .Resource.CamelcasePlural }}(ctx context.Context,{{if ne .Parent nil}}parentID uuid.UUID,{{end}} pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error) {
	offset := (pageNumber - 1) * pageSize

	items, err := s.dbx.FetchRecent{{ .Resource.CamelcasePlural }}(ctx, dbx.FetchRecent{{ .Resource.CamelcasePlural }}Params{
		PageOffset: offset,
		PageLimit:  pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get recent {{ .Resource.CamelcasePlural }}: %w", err)
	}

	totalCount, err := s.dbx.CountRecent{{ .Resource.CamelcasePlural }}(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch {{ .Resource.CamelcasePlural }} count: %w", err)
	}

	return items, totalCount, nil
}
`

const fetchServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Fetch{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID)(dbx.{{ .Resource.CamelcaseSingular }}, error) {
	item, err := s.dbx.Fetch{{ .Resource.CamelcaseSingular }}ByID(ctx, id)
	if err != nil {
		return dbx.{{ .Resource.CamelcaseSingular }}{}, fmt.Errorf("failed to fetch {{ .Resource.CamelcaseSingular }}: %w", err)
	}

	return item, nil
}
`

const destroyServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Destroy{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID) error {
	folderPath := path.Join(s.storageFolder, "{{ .Resource.UnderscoreSingular }}", id.String())

	if err := os.RemoveAll(folderPath); err != nil {
		return fmt.Errorf("failed to remove {{ .Resource.CamelcaseSingular }} folder: %w", err)
	}

	err := s.dbx.Delete{{ .Resource.CamelcaseSingular }}(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to fetch {{ .Resource.CamelcaseSingular }}: %w", err)
	}

	return nil
}
`

type templateDetails struct {
	filename string
	template string
	input    any
}

//nolint:funlen
func (s *Service) generateServiceMethods(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "service", input.Service.String())

	files := map[string]templateDetails{
		"createServiceMethodTemplate": {
			filename: "create_" + input.Resource.UnderscoreSingular() + ".go",
			template: createServiceMethodTemplate,
			input:    input,
		},
		"recentServiceMethodTemplate": {
			filename: "fetch_recent_" + input.Resource.UnderscorePlural() + ".go",
			template: recentServiceMethodTemplate,
			input:    input,
		},
		"fetchServiceMethod": {
			filename: "fetch_" + input.Resource.UnderscoreSingular() + ".go",
			template: fetchServiceMethodTemplate,
			input:    input,
		},
		"destroyServiceMethod": {
			filename: "destroy_" + input.Resource.UnderscoreSingular() + ".go",
			template: destroyServiceMethodTemplate,
			input:    input,
		},
	}

	if input.SearchField != "" {
		files["searchServiceMethod"] = templateDetails{
			filename: fmt.Sprintf("search_%s.go", input.Resource.UnderscorePlural()),
			template: searchServiceMethodTemplate,
			input:    input,
		}
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if field.Type == FieldTypeAttachment {
				files[fmt.Sprintf("update%sServiceMethod", field.Name.CamelcaseSingular())] = templateDetails{
					filename: fmt.Sprintf("upload_%s_%s.go", field.Resource.UnderscoreSingular(), field.Name.UnderscoreSingular()),
					template: uploadAttachmentServiceMethodTemplate,
					input:    field,
				}
			} else {
				files[fmt.Sprintf("update%sServiceMethod", field.Name.CamelcaseSingular())] = templateDetails{
					filename: fmt.Sprintf("update_%s_%s.go", field.Resource.UnderscoreSingular(), field.Name.UnderscoreSingular()),
					template: updateServiceMethodTemplate,
					input:    field,
				}
			}
		}
	}

	for templateName, f := range files {
		filePath := filepath.Join(folderPath, f.filename)
		if err := s.appendTemplateToFile(
			ctx,
			filePath,
			0,
			"",
			templateName,
			f.template,
			f.input,
		); err != nil {
			return fmt.Errorf("failed to append service method: %w", err)
		}

		if err := s.runCommand(folderPath, "goimports", "-w", f.filename); err != nil {
			return fmt.Errorf("failed running goimports: %w", err)
		}
	}

	return nil
}

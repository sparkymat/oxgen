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

const searchServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) Search{{ .Resource.CamelcasePlural }}(ctx context.Context, query string, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error) {
	offset := (pageNumber - 1) * pageSize

	items, err := s.dbx.Search{{ .Resource.CamelcasePlural }}(ctx, dbx.Search{{ .Resource.CamelcasePlural }}Params{
		Query:      query,
		PageOffset: offset,
		PageLimit:  pageSize,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search {{ .Resource.CamelcasePlural }}: %w", err)
	}

	totalCount, err := s.dbx.CountSearched{{ .Resource.CamelcasePlural }}(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch {{ .Resource.CamelcasePlural }} search count: %w", err)
	}

	return items, totalCount, nil
}
`

const recentServiceMethodTemplate = `
package {{ .Service }}

func (s *Service) FetchRecent{{ .Resource.CamelcasePlural }}(ctx context.Context, pageSize int32, pageNumber int32)([]dbx.{{ .Resource.CamelcaseSingular }}, int64, error) {
	offset := (pageNumber - 1) * pageSize

	items, err := s.dbx.FetchRecent{{ .Resource.CamelcasePlural }}(ctx, dbx.Search{{ .Resource.CamelcasePlural }}Params{
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
			filename = fmt.Sprintf("update_%s_%s.go", field.Resource.UnderscoreSingular(), field.Name.UnderscoreSingular())
			filePath = filepath.Join(folderPath, filename)
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

	if input.SearchField != "" {
		// Search
		filename = fmt.Sprintf("search_%s.go", input.Resource.UnderscorePlural())
		filePath = filepath.Join(folderPath, filename)
		if err := s.appendTemplateToFile(
			ctx,
			filePath,
			0,
			"",
			"search",
			searchServiceMethodTemplate,
			input,
		); err != nil {
			return fmt.Errorf("failed to append search service method: %w", err)
		}

		if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
			return fmt.Errorf("failed running goimports: %w", err)
		}
	}

	// Fetch recent
	filename = fmt.Sprintf("fetch_recent_%s.go", input.Resource.UnderscorePlural())
	filePath = filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"recent",
		recentServiceMethodTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to append fetch recents service method: %w", err)
	}

	if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}
	return nil
}

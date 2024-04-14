package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const dbMethodsTemplate = `
  {{if .HasSearch}}
  CountSearched{{ .Resource.CamelcasePlural }}(ctx context.Context, query string) ([]int64, error) 
  {{end}}
  Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params dbx.Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error)
  Delete{{ .Resource.CamelcaseSingular }}(ctx context.Context, id uuid.UUID) error 
  Fetch{{ .Resource.CamelcaseSingular }}ByID(ctx context.Context, id uuid.UUID) (dbx.{{ .Resource.CamelcaseSingular }}, error) 
  Fetch{{ .Resource.CamelcasePlural }}ByIDs(ctx context.Context, ids []uuid.UUID) ([]dbx.{{ .Resource.CamelcaseSingular }}, error) 
  {{if .HasSearch}}
  Search{{ .Resource.CamelcasePlural }}(ctx context.Context, arg dbx.Search{{ .Resource.CamelcasePlural }}Params) ([]dbx.{{ .Resource.CamelcaseSingular }}, error) 
  {{end}}
`

const updateDBMethodTemplate = `
  Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamelcaseSingular }}(ctx context.Context, arg dbx.Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error)
`

func (s *Service) appendDBMethodsToIface(ctx context.Context, workspaceFolder string, name string, fields []Field, searchField string) error {
	input := TemplateInputFromNameAndFields(name, fields, searchField)

	ifaceFilePath := filepath.Join(workspaceFolder, "internal", "service", "database_iface.go")

	if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "dbMethods", dbMethodsTemplate, input); err != nil {
		return err
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "updateDbMethod", updateDBMethodTemplate, field); err != nil {
				return fmt.Errorf("failed to generate update %s SQL method: %w", field.Field.String(), err)
			}
		}
	}

	return nil
}

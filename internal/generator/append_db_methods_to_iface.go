package generator

import (
	"context"
	"path/filepath"
)

const templatedbMethods = `
Create{{ .Resource.CamelcaseSingular }}(ctx context.Context, params dbx.Create{{ .Resource.CamelcaseSingular }}Params) (dbx.{{ .Resource.CamelcaseSingular }}, error)
`

func (s *Service) appendDBMethodsToIface(ctx context.Context, workspaceFolder string, name string, fields []Field, searchField string) error {
	input := TemplateInputFromNameAndFields(name, fields, searchField)

	ifaceFilePath := filepath.Join(workspaceFolder, "internal", "service", "database_iface.go")

	if err := s.appendTemplateToFile(ctx, ifaceFilePath, 2, "}", "dbMethods", templatedbMethods, input); err != nil {
		return err
	}

	return nil
}

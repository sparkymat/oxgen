package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const createSQLMethodTemplate = `-- name: Create{{ .ResourceCamelcaseSingular }} :one
INSERT INTO {{ .ResourceUnderscorePlural }}
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}{{ $field.Name }}{{ end }})
VALUES
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}@{{ $field.Name }}::{{ $field.Type }}{{ end }})
RETURNING *;
`

func (s *Service) generateSQLMethods(ctx context.Context, workspaceFolder string, name string, fields []Field) error {
	input := TemplateInputFromNameAndFields(name, fields)

	if err := s.generateSQLMethod(ctx, workspaceFolder, "create", createSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate create SQL method: %w", err)
	}

	return nil
}

func (s *Service) generateSQLMethod(ctx context.Context, workspaceFolder string, templateName string, templateString string, input TemplateInput) error {
	// Create
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	queriesFile, err := os.OpenFile(filepath.Join(workspaceFolder, "internal", "database", "queries.sql"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open queries file: %w", err)
	}

	if err = tmpl.Execute(queriesFile, input); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

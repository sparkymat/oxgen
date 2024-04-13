package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const createSQLMethodTemplate = `
-- name: Create{{ .ResourceCamelcaseSingular }} :one
INSERT INTO {{ .ResourceUnderscorePlural }}
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}{{ $field.Name }}{{ end }})
VALUES
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}@{{ $field.Name }}::{{ $field.Type }}{{ end }})
RETURNING *;
`

const searchSQLMethodTemplate = `
-- name: Search{{ .ResourceCamelcasePlural }} :many
SELECT *
  FROM {{ .ResourceUnderscorePlural }} t
  WHERE t.{{ .SearchField }} ILIKE '%' || @query::text || '%'
  ORDER BY t.{{ .SearchField }} ASC
  LIMIT @page_limit::int
  OFFSET @page_offset::int;
`

const countSearchedSQLMethodTemplate = `
-- name: CountSearched{{ .ResourceCamelcasePlural }} :many
SELECT COUNT(id)
  FROM {{ .ResourceUnderscorePlural }} t
  WHERE t.{{ .SearchField }} ILIKE '%' || @query::text || '%';
`

const fetchByIdSQLMethodTemplate = `
-- name: Fetch{{ .ResourceCamelcaseSingular }}ByID :one
SELECT *
  FROM {{ .ResourceUnderscorePlural }} t
  WHERE id = @id::uuid
  LIMIT 1;
`

const fetchByIdsSQLMethodTemplate = `
-- name: Fetch{{ .ResourceCamelcasePlural }}ByID :many
SELECT *
  FROM {{ .ResourceUnderscorePlural }} t
  WHERE id = ANY(@ids::uuid[]);
`

func (s *Service) generateSQLMethods(ctx context.Context, workspaceFolder string, name string, fields []Field, searchField string) error {
	input := TemplateInputFromNameAndFields(name, fields, searchField)

	if err := s.generateSQLMethod(ctx, workspaceFolder, "create", createSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate create SQL method: %w", err)
	}

	if searchField != "" {
		if err := s.generateSQLMethod(ctx, workspaceFolder, "search", searchSQLMethodTemplate, input); err != nil {
			return fmt.Errorf("failed to generate search SQL method: %w", err)
		}

		if err := s.generateSQLMethod(ctx, workspaceFolder, "countSearched", countSearchedSQLMethodTemplate, input); err != nil {
			return fmt.Errorf("failed to generate count searched SQL method: %w", err)
		}
	}

	if err := s.generateSQLMethod(ctx, workspaceFolder, "fetchById", fetchByIdSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchById SQL method: %w", err)
	}

	if err := s.generateSQLMethod(ctx, workspaceFolder, "fetchByIds", fetchByIdsSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchByIds SQL method: %w", err)
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
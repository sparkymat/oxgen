package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

const createSQLMethodTemplate = `
-- name: Create{{ .Resource.CamelcaseSingular }} :one
INSERT INTO {{ .Resource.UnderscorePlural }}
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}{{ $field.Field.String }}{{ end }})
VALUES
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}@{{ $field.Field.String }}::{{ $field.Type }}{{ end }})
RETURNING *;
`

const searchSQLMethodTemplate = `
-- name: Search{{ .Resource.CamelcasePlural }} :many
SELECT *
  FROM {{ .Resource.UnderscorePlural }} t
  WHERE t.{{ .SearchField }} ILIKE '%' || @query::text || '%'
  ORDER BY t.{{ .SearchField }} ASC
  LIMIT @page_limit::int
  OFFSET @page_offset::int;
`

const countSearchedSQLMethodTemplate = `
-- name: CountSearched{{ .Resource.CamelcasePlural }} :many
SELECT COUNT(id)
  FROM {{ .Resource.UnderscorePlural }} t
  WHERE t.{{ .SearchField }} ILIKE '%' || @query::text || '%';
`

const fetchByIDSQLMethodTemplate = `
-- name: Fetch{{ .Resource.CamelcaseSingular }}ByID :one
SELECT *
  FROM {{ .Resource.UnderscorePlural }} t
  WHERE id = @id::uuid
  LIMIT 1;
`

const fetchByIDsSQLMethodTemplate = `
-- name: Fetch{{ .Resource.CamelcasePlural }}ByIDs :many
SELECT *
  FROM {{ .Resource.UnderscorePlural }} t
  WHERE id = ANY(@ids::uuid[]);
`

const deleteSQLMethodTemplate = `
-- name: Delete{{ .Resource.CamelcaseSingular }} :exec
DELETE FROM {{ .Resource.UnderscorePlural }} t
  WHERE id = @id::uuid;
`

const updateSQLMethodTemplate = `
-- name: Update{{ .Resource.CamelcaseSingular }}{{ .Field.CamecaseSingular }} :one
UPDATE {{ .Resource.UnderscorePlural }} t
SET {{ .Field.String }} = @{{ .Field.String }}::{{ .Type }}
WHERE id = @id::uuid;
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

	if err := s.generateSQLMethod(ctx, workspaceFolder, "fetchById", fetchByIDSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchById SQL method: %w", err)
	}

	if err := s.generateSQLMethod(ctx, workspaceFolder, "fetchByIds", fetchByIDsSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchByIds SQL method: %w", err)
	}

	if err := s.generateSQLMethod(ctx, workspaceFolder, "delete", deleteSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate delete SQL method: %w", err)
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if err := s.generateSQLMethod(ctx, workspaceFolder, "update", updateSQLMethodTemplate, input); err != nil {
				return fmt.Errorf("failed to generate update %s SQL method: %w", field.Field.String(), err)
			}
		}
	}

	return nil
}

func (*Service) generateSQLMethod(_ context.Context, workspaceFolder string, templateName string, templateString string, input TemplateInput) error {
	// Create
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	//nolint:gosec
	queriesFile, err := os.OpenFile(filepath.Join(workspaceFolder, "internal", "database", "queries.sql"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open queries file: %w", err)
	}

	if err = tmpl.Execute(queriesFile, input); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const createSQLMethodTemplate = `
-- name: Create{{ .Resource.CamelcaseSingular }} :one
INSERT INTO {{ .Resource.UnderscorePlural }}
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}{{ $field.Name.String }}{{ end }})
VALUES
({{ range $i, $field := .Fields }}{{ if $i }}, {{ end }}@{{ $field.Name.String }}::{{ $field.Type }}{{ end }})
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
-- name: CountSearched{{ .Resource.CamelcasePlural }} :one
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
-- name: Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }} :one
UPDATE {{ .Resource.UnderscorePlural }} t
SET {{ .UpdateAssignParamGoFragment }}
WHERE id = @id::uuid
RETURNING *;
`

func (s *Service) generateSQLMethods(ctx context.Context, input Input) error {
	queriesFilePath := filepath.Join(input.WorkspaceFolder, "internal", "database", "queries.sql")

	if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "create", createSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate create SQL method: %w", err)
	}

	if input.SearchField != "" {
		if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "search", searchSQLMethodTemplate, input); err != nil {
			return fmt.Errorf("failed to generate search SQL method: %w", err)
		}

		if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "countSearched", countSearchedSQLMethodTemplate, input); err != nil {
			return fmt.Errorf("failed to generate count searched SQL method: %w", err)
		}
	}

	if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "fetchById", fetchByIDSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchById SQL method: %w", err)
	}

	if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "fetchByIds", fetchByIDsSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate fetchByIds SQL method: %w", err)
	}

	if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "delete", deleteSQLMethodTemplate, input); err != nil {
		return fmt.Errorf("failed to generate delete SQL method: %w", err)
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if err := s.appendTemplateToFile(ctx, queriesFilePath, 0, "", "update", updateSQLMethodTemplate, field); err != nil {
				return fmt.Errorf("failed to generate update %s SQL method: %w", field.Name.String(), err)
			}
		}
	}

	return nil
}

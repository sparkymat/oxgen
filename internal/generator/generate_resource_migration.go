package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const upTemplate = `CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE {{ .Resource.UnderscorePlural }} (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
{{range .Fields}}{{ .CreateSQLFragment }},
{{end}}  created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER {{ .Resource.UnderscorePlural }}_updated_at
  BEFORE UPDATE
  ON {{ .Resource.UnderscorePlural }}
  FOR EACH ROW
    EXECUTE FUNCTION moddatetime(updated_at);
`

const downTemplate = `DROP TABLE {{ .Resource.UnderscorePlural }};
`

func (s *Service) generateResourceMigration(_ context.Context, input Input) error {
	if err := s.ensureFolderExists(filepath.Join(input.WorkspaceFolder, "migrations")); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")

	// up
	upTmpl, err := template.New("up").Parse(upTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse up template: %w", err)
	}

	//nolint:gosec
	upFile, err := os.Create(
		filepath.Join(
			input.WorkspaceFolder,
			"migrations",
			fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, input.Resource.UnderscorePlural()),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create up file: %w", err)
	}

	if err = upTmpl.Execute(upFile, input); err != nil {
		return fmt.Errorf("failed to execute up template: %w", err)
	}

	// down
	downTmpl, err := template.New("down").Parse(downTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse down template: %w", err)
	}

	//nolint:gosec
	downFile, err := os.Create(
		filepath.Join(
			input.WorkspaceFolder,
			"migrations",
			fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, input.Resource.UnderscorePlural()),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create up file: %w", err)
	}

	if err = downTmpl.Execute(downFile, input); err != nil {
		return fmt.Errorf("failed to execute up template: %w", err)
	}

	return nil
}

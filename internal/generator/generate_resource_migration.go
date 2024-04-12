package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/samber/lo"
)

const upTemplate = `CREATE EXTENSION IF NOT EXISTS moddatetime;

CREATE TABLE {{.ResourceUnderscorePlural}} (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
{{range .Fields}}  {{.Name}} {{.Type}} {{.Modifiers}}{{.Default}},
{{end}}  created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER {{.ResourceUnderscorePlural}}_updated_at
  BEFORE UPDATE
  ON {{.ResourceUnderscorePlural}}
  FOR EACH ROW
    EXECUTE FUNCTION moddatetime(updated_at);
`

const downTemplate = `DROP TABLE {{.ResourceUnderscorePlural}};
`

func (s *Service) generateResourceMigration(_ context.Context, name string, fields []Field) error {
	if err := s.ensureFolderExists("migrations"); err != nil {
		return err
	}

	input := MigrationTemplateInput{
		ResourceUnderscorePlural: pluralize.NewClient().Plural(strcase.ToSnake(name)),
		Fields: lo.Map(fields, func(f Field, _ int) MigrationTemplateInputField {
			return f.MigrationTemplateInputField()
		}),
	}

	timestamp := time.Now().Format("20060102150405")

	// up
	upTmpl, err := template.New("up").Parse(upTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse up template: %w", err)
	}

	upFile, err := os.Create(filepath.Join("migrations", fmt.Sprintf("%s_create_%s_table.up.sql", timestamp, input.ResourceUnderscorePlural)))
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

	downFile, err := os.Create(filepath.Join("migrations", fmt.Sprintf("%s_create_%s_table.down.sql", timestamp, input.ResourceUnderscorePlural)))
	if err != nil {
		return fmt.Errorf("failed to create up file: %w", err)
	}

	if err = downTmpl.Execute(downFile, input); err != nil {
		return fmt.Errorf("failed to execute up template: %w", err)
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const frontendModelTemplate = `
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';

dayjs.extend(utc);
{{range $f := .Fields}}{{ if eq $f.Type "enum" }}
{{ $f.EnumTypesFrontendModel }}
{{end}}{{end}}

export class {{  .Resource.CamelcaseSingular }} {
  public id: string;

  public createdAt: dayjs.Dayjs;

  public updatedAt: dayjs.Dayjs;

  {{range .Fields }}{{ .FrontendModelDeclaration }}
  {{end}}

  constructor(json: any) {
    if (!json) {
      return;
    }

    this.id = json.id;
    this.createdAt = dayjs.utc(json.createdAt);
    this.updatedAt = dayjs.utc(json.updatedAt);

    {{range .Fields }}{{ .FrontendModelAssignment }}
    {{end}}
  }
}
`

func (s *Service) generateFrontendModel(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "frontend", "src", "models")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure frontend models folder exists: %w", err)
	}

	filename := input.Resource.CamelcaseSingular() + ".ts"

	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"frontendModel",
		frontendModelTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to generate frontend model: %w", err)
	}

	return nil
}

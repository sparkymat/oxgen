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

export class {{  .Resource.CamelcaseSingular }} {
  public id: string;

  public createdAt: string;

  public updatedAt: string;

  {{range .Fields }}{{ .FrontendModelDeclaration }}
  {{end}}

  construction(json: any) {
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

	filename := input.Resource.UnderscoreSingular() + ".ts"

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
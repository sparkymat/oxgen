package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const routeMethodsTemplate = `
  apiGroup.POST("/{{ .Resource.UnderscorePlural }}", api.{{ .Resource.CamelcasePlural }}Create(services))
  {{if .HasSearch}}apiGroup.GET("/{{ .Resource.UnderscorePlural }}/search", api.{{ .Resource.CamelcasePlural }}Search(services))
  {{end}}apiGroup.GET("/{{ .Resource.UnderscorePlural }}/recent", api.{{ .Resource.CamelcasePlural }}FetchRecent(services))
  apiGroup.GET("/{{ .Resource.UnderscorePlural }}/:id", api.{{ .Resource.CamelcasePlural }}Show(services))
  apiGroup.DELETE("/{{ .Resource.UnderscorePlural }}/:id", api.{{ .Resource.CamelcasePlural }}Destroy(services))
`

const updateRouteMethodTemplate = `
  apiGroup.PATCH("/{{ .Resource.UnderscorePlural }}/:id/update_{{ .Name.UnderscoreSingular }}", api.{{ .Resource.CamelcasePlural }}Update{{ .Name.CamelcaseSingular }}(services))
`

const uploadRouteMethodTemplate = `
  apiGroup.PATCH("/{{ .Resource.UnderscorePlural }}/:id/upload_{{ .Name.UnderscoreSingular }}", api.{{ .Resource.CamelcasePlural }}Upload{{ .Name.CamelcaseSingular }}(services))
`

func (s *Service) appendRoutes(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "route")
	filename := "api.go"

	filePath := filepath.Join(folderPath, filename)

	if err := s.appendTemplateToFile(ctx, filePath, 2, "}", "routeMethods", routeMethodsTemplate, input); err != nil {
		return fmt.Errorf("failed to generate route methods: %w", err)
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if field.Type == FieldTypeAttachment {
				if err := s.appendTemplateToFile(ctx, filePath, 2, "}", "uploadRouteMethod", uploadRouteMethodTemplate, field); err != nil {
					return fmt.Errorf("failed to generate update %s route: %w", field.Name.String(), err)
				}
			} else {
				if err := s.appendTemplateToFile(ctx, filePath, 2, "}", "updateRouteMethod", updateRouteMethodTemplate, field); err != nil {
					return fmt.Errorf("failed to generate update %s route: %w", field.Name.String(), err)
				}
			}
		}
	}

	if err := s.runCommand(folderPath, "goimports", "-w", filename); err != nil {
		return fmt.Errorf("failed running goimports: %w", err)
	}

	return nil
}

package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const createHandlerMethodTemplate = `
package api

type {{ .Resource.CamelcasePlural }}CreateRequest struct {
{{range .Fields }}{{if .Initial}}{{ .CreateRequestGoFragment }}{{end}}
{{end}}
}

func {{ .Resource.CamelcasePlural }}Create(s internal.Services) echo.HandlerFunc {
  return wrapWithAuth(func(c echo.Context, _ dbx.User) error {
    var request {{ .Resource.CamelcasePlural }}CreateRequest
    if err := c.Bind(&request); err != nil {
      return renderError(c, http.StatusBadRequest, "invalid request", err)
    }

    item, err := s.{{ .Service.Capitalize }}.Create{{ .Resource.CamelcaseSingular }}(
      c.Request().Context(), 
    )
    if err != nil {
      return renderError(c, http.StatusInternalServerError, "could not create {{ .Resource.CamelcaseSingular }}", err)
    }

    presented{{ .Resource.CamelcaseSingular }} := presenter.{{ .Resource.CamelcaseSingular }}FromModel(item)

    return c.JSON(http.StatusCreated, presented{{ .Resource.CamelcaseSingular }})
  })
}
`

const updateHandlerMethodTemplate = `
package api

type {{ .Resource.CamelcasePlural }}Update{{ .Name.CamelcaseSingular }}Request struct {
  {{ .CreateRequestGoFragment }}
}

func {{ .Resource.CamelcasePlural }}Update{{ .Name.CamelcaseSingular }}(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
		var request {{ .Resource.CamelcasePlural }}Update{{ .Name.CamelcaseSingular }}Request
		if err := c.Bind(&request); err != nil {
			return renderError(c, http.StatusBadRequest, "invalid request", err)
		}

		item, err := s.{{ .Service.Capitalize }}.Update{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(c.Request().Context(), id, request.{{ .Name.CamelcaseSingular }})
    if err != nil {
      return renderError(c, http.StatusInternalServerError, "could not update {{ .Resource.CamelcaseSingular }} {{ .Name.CamelcaseSingular }}", err)
		}

    presented{{ .Resource.CamelcaseSingular }} := presenter.{{ .Resource.CamelcaseSingular }}FromModel(item)

    return c.JSON(http.StatusCreated, presented{{ .Resource.CamelcaseSingular }})
  })
}
`

const uploadAttachmentHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}Upload{{ .Name.CamelcaseSingular }}(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
		fileHeader, err := c.FormFile("{{ .Name.UnderscoreSingular }}_file")
		if err != nil {
			return renderError(c, http.StatusBadRequest, "missing file", err)
		}

		file, err := fileHeader.Open()
		if err != nil {
			return renderError(c, http.StatusBadRequest, "failed to open file", err)
		}
		defer file.Close()

		item, err := s.{{ .Service.Capitalize }}.Upload{{ .Resource.CamelcaseSingular }}{{ .Name.CamelcaseSingular }}(
			c.Request().Context(),
			id,
			fileHeader.Filename,
			file,
		)
    if err != nil {
			return renderError(c, http.StatusInternalServerError, "failed to upload {{ .Name.UnderscoreSingular }}", err)
		}

    presented{{ .Resource.CamelcaseSingular }} := presenter.{{ .Resource.CamelcaseSingular }}FromModel(item)

    return c.JSON(http.StatusCreated, presented{{ .Resource.CamelcaseSingular }})
  })
}
`

const searchHandlerMethodTemplate = `
package api

type {{ .Resource.CamelcasePlural }}SearchResponse struct {
  Items []presenter.{{ .Resource.CamelcaseSingular }} ` + "`json:\"items\"`" + `
  TotalCount int ` + "`json:\"totalCount\"`" + `
  PageSize int ` + "`json:\"pageSize\"`" + `
  PageNumber int ` + "`json:\"pageNumber\"`" + `
}

func {{ .Resource.CamelcasePlural }}Search(s internal.Services) echo.HandlerFunc {
  return wrapWithAuth(func(c echo.Context, _ dbx.User) error {
  })
}
`

const recentHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}FetchRecent(s internal.Services) echo.HandlerFunc {
  return wrapWithAuth(func(c echo.Context, _ dbx.User) error {
  })
}
`

const fetchHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}Show(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
  })
}
`

const destroyHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}Destroy(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
  })
}
`

func (s *Service) generateHandlerMethods(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "internal", "handler", "api")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure service folder exists: %w", err)
	}

	files := map[string]templateDetails{
		"createHandlerMethodTemplate": {
			filename: input.Resource.UnderscorePlural() + "_create.go",
			template: createHandlerMethodTemplate,
			input:    input,
		},
		"recentHandlerMethodTemplate": {
			filename: input.Resource.UnderscorePlural() + "_fetch_recent.go",
			template: recentHandlerMethodTemplate,
			input:    input,
		},
		"fetchHandlerMethod": {
			filename: input.Resource.UnderscorePlural() + "_show.go",
			template: fetchHandlerMethodTemplate,
			input:    input,
		},
		"destroyHandlerMethod": {
			filename: input.Resource.UnderscorePlural() + "_destroy.go",
			template: destroyHandlerMethodTemplate,
			input:    input,
		},
	}

	if input.SearchField != "" {
		files["searchHandlerMethod"] = templateDetails{
			filename: input.Resource.UnderscorePlural() + "_search.go",
			template: searchHandlerMethodTemplate,
			input:    input,
		}
	}

	for _, field := range input.Fields {
		if field.Updateable {
			if field.Type == FieldTypeAttachment {
				files[fmt.Sprintf("update%sHandlerMethod", field.Name.CamelcaseSingular())] = templateDetails{
					filename: fmt.Sprintf("%s_upload_%s.go", field.Resource.UnderscorePlural(), field.Name.UnderscoreSingular()),
					template: uploadAttachmentHandlerMethodTemplate,
					input:    field,
				}
			} else {
				files[fmt.Sprintf("update%sHandlerMethod", field.Name.CamelcaseSingular())] = templateDetails{
					filename: fmt.Sprintf("%s_update_%s.go", field.Resource.UnderscorePlural(), field.Name.UnderscoreSingular()),
					template: updateHandlerMethodTemplate,
					input:    field,
				}
			}
		}
	}

	for templateName, f := range files {
		filePath := filepath.Join(folderPath, f.filename)
		if err := s.appendTemplateToFile(
			ctx,
			filePath,
			0,
			"",
			templateName,
			f.template,
			f.input,
		); err != nil {
			return fmt.Errorf("failed to generate handler method: %w", err)
		}

		if err := s.runCommand(folderPath, "goimports", "-w", f.filename); err != nil {
			return fmt.Errorf("failed running goimports: %w", err)
		}
	}

	return nil
}

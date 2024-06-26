//nolint:lll,revive
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

    input := {{ .Service.String }}.Create{{ .Resource.CamelcaseSingular }}Params{}

    {{range .Fields }}{{if .Initial}}{{ .CreateHandlerAssignParamsGoFragment }}
    {{end}}{{end}}

    item, err := s.{{ .Service.Capitalize }}.Create{{ .Resource.CamelcaseSingular }}(
      c.Request().Context(), 
      input,
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
{{if eq .Parent nil}}  return wrapWithAuth(func(c echo.Context, _ dbx.User) error {
{{else}}  return wrapWithAuthForChild(func(c echo.Context, _ dbx.User, parentID uuid.UUID) error {
{{end}}
		pageSize, pageNumber, err := parsePaginationParams(c)
		if err != nil {
			return renderError(c, http.StatusBadRequest, "invalid pagination params", err)
		}

		query := c.QueryParam("query")

		items, totalCount, err := s.{{ .Service.Capitalize }}.Search{{ .Resource.CamelcasePlural }}(c.Request().Context(),{{if ne .Parent nil}} parentID,{{end}} query, pageSize, pageNumber)
		if err != nil {
			return renderError(c, http.StatusInternalServerError, "failed to search {{ .Resource.LowerCamelcasePlural }}", err)
		}

		presentedItems := lo.Map(items, func(i dbx.{{ .Resource.CamelcaseSingular }}, _ int) presenter.{{ .Resource.CamelcaseSingular }} {
			return presenter.{{ .Resource.CamelcaseSingular }}FromModel(i)
		})

		response := {{ .Resource.CamelcasePlural }}SearchResponse{
			Items:      presentedItems,
			PageSize:   int(pageSize),
			PageNumber: int(pageNumber),
			TotalCount: int(totalCount),
		}

		return c.JSON(http.StatusOK, response)
  })
}
`

const recentHandlerMethodTemplate = `
package api

type {{ .Resource.CamelcasePlural }}FetchRecentResponse struct {
  Items []presenter.{{ .Resource.CamelcaseSingular }} ` + "`json:\"items\"`" + `
  TotalCount int ` + "`json:\"totalCount\"`" + `
  PageSize int ` + "`json:\"pageSize\"`" + `
  PageNumber int ` + "`json:\"pageNumber\"`" + `
}

func {{ .Resource.CamelcasePlural }}FetchRecent(s internal.Services) echo.HandlerFunc {
{{if eq .Parent nil}}  return wrapWithAuth(func(c echo.Context, _ dbx.User) error {
{{else}}  return wrapWithAuthForChild(func(c echo.Context, _ dbx.User, parentID uuid.UUID) error {
{{end}}
		pageSize, pageNumber, err := parsePaginationParams(c)
		if err != nil {
			return renderError(c, http.StatusBadRequest, "invalid pagination params", err)
		}

		items, totalCount, err := s.{{ .Service.Capitalize }}.FetchRecent{{ .Resource.CamelcasePlural }}(c.Request().Context(),{{if ne .Parent nil}} parentID,{{end}} pageSize, pageNumber)
		if err != nil {
			return renderError(c, http.StatusInternalServerError, "failed to fetch recent {{ .Resource.LowerCamelcasePlural }}", err)
		}

		presentedItems := lo.Map(items, func(i dbx.{{ .Resource.CamelcaseSingular }}, _ int) presenter.{{ .Resource.CamelcaseSingular }} {
			return presenter.{{ .Resource.CamelcaseSingular }}FromModel(i)
		})

		response := {{ .Resource.CamelcasePlural }}FetchRecentResponse{
			Items:      presentedItems,
			PageSize:   int(pageSize),
			PageNumber: int(pageNumber),
			TotalCount: int(totalCount),
		}

		return c.JSON(http.StatusOK, response)
  })
}
`

const fetchHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}Show(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
    item, err := s.{{ .Service.Capitalize }}.Fetch{{ .Resource.CamelcaseSingular }}(
			c.Request().Context(),
			id,
		)
		if err != nil {
			return renderError(c, http.StatusInternalServerError, "failed to fetch {{ .Resource.LowerCamelcaseSingular }}", err)
		}

		presentedItem := presenter.{{ .Resource.CamelcaseSingular }}FromModel(item)

		return c.JSON(http.StatusOK, presentedItem)
  })
}
`

const destroyHandlerMethodTemplate = `
package api

func {{ .Resource.CamelcasePlural }}Destroy(s internal.Services) echo.HandlerFunc {
  return wrapWithAuthForMember(func(c echo.Context, _ dbx.User, id uuid.UUID) error {
		if err := s.{{ .Service.Capitalize }}.Destroy{{ .Resource.CamelcaseSingular }}(c.Request().Context(), id); err != nil {
			return renderError(c, http.StatusInternalServerError, "failed to destroy {{ .Resource.LowerCamelcaseSingular }}", err)
		}

		return c.NoContent(http.StatusOK)
  })
}
`

//nolint:funlen
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

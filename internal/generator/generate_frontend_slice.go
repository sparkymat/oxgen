package generator

import (
	"context"
	"fmt"
	"path/filepath"
)

const frontendSliceTemplate = `
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { {{ .Resource.CamelcaseSingular }} } from '../models/{{ .Resource.CamelcaseSingular }}';

export interface ListResponse {
  items: {{ .Resource.CamelcaseSingular }}[];
  totalCount: number;
  pageNumber: number;
  pageSize: number;
}

export interface FetchRecentRequest {
  pageSize: number;
  pageNumber: number;
}

{{if .HasSearch }}
export interface SearchRequest {
  query: string;
  pageSize: number;
  pageNumber: number;
}
{{end}}

export interface CreateRequest {
{{range .Fields }}{{if .Initial}}{{ .FrontendModelDeclaration }}
{{end}}{{end}}
}

{{range .Fields}}{{if .Updateable}}{{if eq .Type "attachment"}}export interface Upload{{ .Name.CamelcaseSingular }}Request {
  id: string;
  formData: FormData;
}

{{else}}export interface Update{{ .Name.CamelcaseSingular }}Request {
  id: string;
{{ .FrontendModelDeclaration }}
}

{{end}}{{end}}{{end}}

export const api = createApi({
  reducerPath: '{{ .Resource.LowerCamelcasePlural }}',
  baseQuery: fetchBaseQuery({ baseUrl: '/api' }),
  endpoints: builder => ({
    recent: builder.query<ListResponse, FetchRecentRequest>({
      query: ({pageSize, pageNumber}) => ` + "`{{ .Resource.UnderscorePlural }}?pageSize=${pageSize}&pageNumber=${pageNumber}`" + `,
      providesTags: [{ type: '{{ .Resource.CamelcasePlural }}', id: 'LIST' }],
    }),
    {{if .HasSearch }}search: builder.query<ListResponse, SearchRequest>({
      query: ({pageSize, pageNumber}) => ` + "`{{ .Resource.UnderscorePlural }}?query=${encodeURIComponent(query)}&pageSize=${pageSize}&pageNumber=${pageNumber}`" + `,
      providesTags: [{ type: '{{ .Resource.CamelcasePlural }}', id: 'SEARCH' }],
    }),{{end}}
    show: builder.query<{{ .Resource.CamelcaseSingular }}, string>({
      query: id => "{{ .Resource.UnderscorePlural }}/${id}",
      providesTags: (_result, _error, arg) => [{ type: '{{ .Resource.CamelcaseSingular }}', id: arg }],
    }),
    create: builder.mutation<{{ .Resource.CamelcaseSingular }}, CreateRequest>({
      query: body => ({
        url: '{{ .Resource.UnderscorePlural }}',
        method: 'POST',
        body,
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content,
        },
      }),
      providesTags: (result, _error, _arg) => [
        { type: '{{ .Resource.CamelcaseSingular }}', id: result.id },
      ],
      invalidatesTags: [{ type: '{{ .Resource.CamelcaseSingular }}', id: 'LIST' }],
    }),
    destroy: builder.mutation<void,string>({
      query: id => ({
        url: '{{ .Resource.UnderscorePlural }}/${id}',
        method: 'DELETE',
        headers: {
          'X-CSRF-Token': (
            document.querySelector('meta[name="csrf-token"]') as any
          ).content,
        },
      }),
      invalidatesTags: (_result, _error, arg) => [{ type: '{{ .Resource.CamelcaseSingular }}', id: arg }],
    }),
    {{range .Fields}}{{if .Updateable}}
      {{if eq .Type "attachment"}}upload{{ .Name.CamelcaseSingular }}: builder.mutation<{{ .Resource.CamelcaseSingular }}, Upload{{ .Name.CamelcaseSingular }}Request>({
        query: ({id, formData }) => ({
          url: ` + "`" + `{{ .Resource.UnderscorePlural }}/${id}/upload_{{ .Resource.UnderscoreSingular }}` + "`" + `,
          method: 'PATCH',
          body: formData,
          headers: {
            'X-CSRF-Token': (
              document.querySelector('meta[name="csrf-token"]') as any
            ).content,
          },
        }),
        providesTags: (result, _error, _arg) => [
          { type: '{{ .Resource.CamelcaseSingular }}', id: result.id },
        ],
        invalidatesTags: [{ type: '{{ .Resource.CamelcaseSingular }}', id: 'LIST' }, { type: '{{ .Resource.CamelcaseSingular }}', id: 'SEARCH' }],
      }),
      {{else}}update{{ .Name.CamelcaseSingular }}: builder.mutation<{{ .Resource.CamelcaseSingular }}, Update{{ .Name.CamelcaseSingular }}Request>({
        query: ({id, {{ .Name.LowerCamelcaseSingular }} }) => ({
          url: ` + "`" + `{{ .Resource.UnderscorePlural }}/${id}/upload_{{ .Resource.UnderscoreSingular }}` + "`" + `,
          method: 'PATCH',
          body: { {{ .Name.LowerCamelcaseSingular }} },
          headers: {
            'X-CSRF-Token': (
              document.querySelector('meta[name="csrf-token"]') as any
            ).content,
          },
        }),
        providesTags: (result, _error, _arg) => [
          { type: '{{ .Resource.CamelcaseSingular }}', id: result.id },
        ],
        invalidatesTags: [{ type: '{{ .Resource.CamelcaseSingular }}', id: 'LIST' }, { type: '{{ .Resource.CamelcaseSingular }}', id: 'SEARCH' }],
      }),{{end}}{{end}}{{end}}
  }),
});

export const {
  useRecentQuery,
  {{if .HasSearch}}useSearchQuery,
  {{end}}useCreateMutation,
  useShowQuery,
  useDestroyMutation
} = api;

`

func (s *Service) generateFrontendSlice(ctx context.Context, input Input) error {
	folderPath := filepath.Join(input.WorkspaceFolder, "frontend", "src", "slices")

	if err := s.ensureFolderExists(folderPath); err != nil {
		return fmt.Errorf("failed to ensure frontend slices folder exists: %w", err)
	}

	filename := input.Resource.CamelcaseSingular() + ".ts"

	filePath := filepath.Join(folderPath, filename)
	if err := s.appendTemplateToFile(
		ctx,
		filePath,
		0,
		"",
		"frontendSlice",
		frontendSliceTemplate,
		input,
	); err != nil {
		return fmt.Errorf("failed to generate frontend slice: %w", err)
	}

	return nil
}

package generator

import "context"

const createServiceMethodTemplate = `
package {{}}

type Create{{}}Params struct {
}

func (s *Service) Create{{}}(ctx context.Context, params dbx.Create{{}}Params) (dbx.{{}}, error) {
  return s.db.Create{{}}(ctx, params)
}
`

func (s *Service) addServiceMethods(ctx context.Context, input GenerateInput) error {
	return nil
}

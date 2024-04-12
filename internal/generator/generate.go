package generator

import (
	"context"
	"errors"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

var ErrInvalidResourceName = errors.New("invalid resource name")

func (s *Service) Generate(ctx context.Context, name string) error {
	if err := ensureValidResourceName(name); err != nil {
		return err
	}

	// migration
	if err := s.generateResourceMigration(ctx, name); err != nil {
		return err
	}

	// run migration, dump schema and generate models

	// add sql methods

	// run sqlc gen

	// copy new methods to database_iface

	// add new methods to service

	// add new methods to handler

	// add new methods to routes

	return nil
}

func ensureValidResourceName(name string) error {
	if name == "" {
		return ErrInvalidResourceName
	}

	if !pluralize.NewClient().IsSingular(name) || strcase.ToCamel(name) != name {
		return ErrInvalidResourceName
	}

	return nil
}

package generator

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

var ErrInvalidResourceName = errors.New("invalid resource name")

func (*Service) CheckValidProject(_ context.Context, workspaceFolder string) error {
	// check if the workspace folder exists
	if info, err := os.Stat(workspaceFolder); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("workspace folder %s does not exist: %w", workspaceFolder, err)
	}

	// Ensure Makefile
	makeFilePath := filepath.Join(workspaceFolder, "Makefile")
	if _, err := os.Stat(makeFilePath); os.IsNotExist(err) {
		return fmt.Errorf("Makefile not found in workspace folder %s: %w", workspaceFolder, err) //nolint:stylecheck
	}

	return nil
}

func (s *Service) Generate(ctx context.Context, workspaceFolder string, name string, fieldStrings []string) error {
	if err := ensureValidResourceName(name); err != nil {
		return err
	}

	fields := []Field{}

	for _, fieldString := range fieldStrings {
		field, err := ParseField(fieldString)
		if err != nil {
			return fmt.Errorf("failed parsing field %s: %w", fieldString, err)
		}

		fields = append(fields, field)
	}

	// migration
	if err := s.generateResourceMigration(ctx, workspaceFolder, name, fields); err != nil {
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

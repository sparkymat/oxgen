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

type GenerateInput struct {
	WorkspaceFolder string
	Service         string
	Name            string
	FieldStrings    []string
	SearchField     string
}

func (*Service) CheckValidProject(_ context.Context, workspaceFolder string) error {
	// check if the workspace folder exists
	if info, err := os.Stat(workspaceFolder); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("workspace folder %s does not exist: %w", workspaceFolder, err)
	}

	// Ensure Makefile
	makeFilePath := filepath.Join(workspaceFolder, "Makefile")
	if _, err := os.Stat(makeFilePath); os.IsNotExist(err) {
		return fmt.Errorf("makefile not found in workspace folder %s: %w", workspaceFolder, err)
	}

	return nil
}

func (s *Service) Generate(
	ctx context.Context,
	input GenerateInput,
) error {
	if err := ensureValidResourceName(input.Name); err != nil {
		return err
	}

	fields := []Field{}

	for _, fieldString := range input.FieldStrings {
		field, err := ParseField(input.Name, fieldString)
		if err != nil {
			return fmt.Errorf("failed parsing field %s: %w", fieldString, err)
		}

		fields = append(fields, field)
	}

	// migration
	if err := s.generateResourceMigration(ctx, input.WorkspaceFolder, input.Name, fields, input.SearchField); err != nil {
		return fmt.Errorf("failed generating resource migration: %w", err)
	}

	// run migration, dump schema and generate models
	if err := s.runCommand(input.WorkspaceFolder, "make", "db-migrate"); err != nil {
		return fmt.Errorf("failed running make db-migrate: %w", err)
	}

	if err := s.runCommand(input.WorkspaceFolder, "make", "db-schema-dump"); err != nil {
		return fmt.Errorf("failed running make db-schema-dump: %w", err)
	}

	// add sql methods
	if err := s.generateSQLMethods(ctx, input.WorkspaceFolder, input.Name, fields, input.SearchField); err != nil {
		return fmt.Errorf("failed generating sql methods: %w", err)
	}

	// run sqlc gen
	if err := s.runCommand(input.WorkspaceFolder, "make", "sqlc-gen"); err != nil {
		return fmt.Errorf("failed running make sqlc-gen: %w", err)
	}

	// copy new methods to database_iface
	if err := s.appendDBMethodsToIface(ctx, input.WorkspaceFolder, input.Name, fields, input.SearchField); err != nil {
		return fmt.Errorf("failed appending new methods to database_iface.go: %w", err)
	}

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

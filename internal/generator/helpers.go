// Package generator provides a service to generate code based on the database schema.
package generator

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidPath = errors.New("invalid path")

func (*Service) ensureFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0o755) //nolint:gosec,gomnd
		if err != nil {
			return fmt.Errorf("failed to create migrations path: %w", err)
		}
	}

	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("migrations path not found: %w", ErrInvalidPath)
	}

	return nil
}

package generator

import (
	"errors"
	"fmt"
	"os"
)

var ErrInvalidPath = errors.New("invalid path")

func (s *Service) ensureFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}

	if info, err := os.Stat(path); os.IsNotExist(err) || !info.IsDir() {
		return fmt.Errorf("migrations path not found: %w", ErrInvalidPath)
	}

	return nil
}

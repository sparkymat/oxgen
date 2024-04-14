// Package generator provides a service to generate code based on the database schema.
package generator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
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

func (*Service) runCommand(workspaceFolder string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = workspaceFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed running %s %s: %w", command, strings.Join(args, " "), err)
	}

	return nil
}

func (*Service) appendTemplateToFile(
	_ context.Context,
	filePath string,
	reverseOffset int,
	suffix string,
	templateName string,
	templateString string,
	input any,
) error {
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var queriesFile *os.File

	if reverseOffset != 0 {
		queriesFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0o644) //nolint:gomnd,gosec
		if err != nil {
			return fmt.Errorf("failed to open queries file: %w", err)
		}

		defer queriesFile.Close()

		if _, err = queriesFile.Seek(int64(-reverseOffset), io.SeekEnd); err != nil {
			return fmt.Errorf("failed to seek to the end of the file: %w", err)
		}
	} else {
		queriesFile, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644) //nolint:gomnd,gosec
		if err != nil {
			return fmt.Errorf("failed to open queries file: %w", err)
		}
	}

	if err = tmpl.Execute(queriesFile, input); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if suffix != "" {
		if _, err := io.WriteString(queriesFile, suffix); err != nil {
			return fmt.Errorf("failed to write suffix: %w", err)
		}
	}

	return nil
}

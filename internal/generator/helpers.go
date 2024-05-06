// Package generator provides a service to generate code based on the database schema.
package generator

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

var (
	ErrInvalidPath    = errors.New("invalid path")
	ErrAnchorNotFound = errors.New("anchor not found")
)

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

func (*Service) ensureFileExists(path string, templateName string, templateString string, templateInput any) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		defer file.Close() //nolint:errcheck

		tmpl, err := template.New(templateName).Parse(templateString)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		if err = tmpl.Execute(file, templateInput); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path not found: %w", ErrInvalidPath)
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

		defer queriesFile.Close() //nolint:errcheck

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

func (*Service) injectTemplateAboveLine(
	filePath string,
	anchorLine string,
	templateName string,
	templateString string,
	input any,
) error {
	tmpl, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	targetFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open target file: %w", err)
	}

	defer targetFile.Close() //nolint:errcheck

	targetLines := []string{}
	scanner := bufio.NewScanner(targetFile)

	anchorIndex := -1
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		targetLines = append(targetLines, line)

		if strings.TrimSpace(line) == anchorLine {
			anchorIndex = lineNumber
		}

		lineNumber += 1
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to read target file: %w", err)
	}

	if anchorIndex == -1 {
		return ErrAnchorNotFound
	}

	templateBuf := &bytes.Buffer{}
	if err = tmpl.Execute(templateBuf, input); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	finalLines := append(targetLines[:anchorIndex+1], targetLines[anchorIndex:]...)
	finalLines[anchorIndex] = templateBuf.String()

	err = writeLinesToFile(finalLines, filePath)
	if err != nil {
		return fmt.Errorf("failed to write target file: %w", err)
	}

	return nil
}

func writeLinesToFile(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

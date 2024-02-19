package template

import (
	"errors"
	"regexp"
)

var (
	ErrConfigNoCommands  = errors.New("no commands")
	ErrConfigInvalidName = errors.New("invalid name")
	ErrConfigNoFiles     = errors.New("no files")
)

type Config struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	Name      string `json:"name"`
	Overwrite bool   `json:"overwrite"`
	Files     []File `json:"files"`
}

type File struct {
	Path     string `json:"path"`
	Template string `json:"template"`
}

func (c Config) Validate() error {
	if len(c.Commands) == 0 {
		return ErrConfigNoCommands
	}

	for _, cmd := range c.Commands {
		if err := cmd.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c Command) Validate() error {
	nameRegex := regexp.MustCompile(`[a-z]{2}[a-z0-9-_]*`)
	if !nameRegex.MatchString(c.Name) {
		return ErrConfigInvalidName
	}

	if len(c.Files) == 0 {
		return ErrConfigNoFiles
	}

	return nil
}

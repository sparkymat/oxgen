package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/sparkymat/oxgen/internal/generator"
	"github.com/sparkymat/oxgen/internal/git"
)

var ErrUncommittedChanges = errors.New("uncommitted changes")

func main() {
	var err error

	configContents, err := os.ReadFile("oxgen.json")
	if err != nil {
		panic(err)
	}

	gitRepo, err := git.New()
	if err != nil {
		panic(err)
	}

	repoClean, err := gitRepo.StatusClean()
	if err != nil {
		panic(err)
	}

	if !repoClean {
		panic(ErrUncommittedChanges)
	}

	var config generator.Config
	if err = json.Unmarshal(configContents, &config); err != nil {
		panic(err)
	}

	s := generator.New(config)

	err = s.Generate(context.Background()) //nolint:wrapcheck
	if err != nil {
		panic(err)
	}
}

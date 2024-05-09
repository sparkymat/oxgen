package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sparkymat/oxgen/internal/generator"
	"github.com/sparkymat/oxgen/internal/git"
	"github.com/spf13/cobra"
)

var ErrUncommittedChanges = errors.New("uncommitted changes")

var skipGitCheck bool //nolint:gochecknoglobals

var service string //nolint:gochecknoglobals

var workspaceFolder string //nolint:gochecknoglobals

var searchField string //nolint:gochecknoglobals

var parent string //nolint:gochecknoglobals

//nolint:gochecknoglobals
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "resource generates a new resource for the project",
	Long:  `resource generates a new resource for the project. `,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("Generating resource")

		if !skipGitCheck {
			log.Info().Msg("Checking for uncommitted changes")
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
		}

		gen := generator.New()

		if err := gen.CheckValidProject(cmd.Context(), workspaceFolder); err != nil {
			panic(err)
		}

		name := args[0]
		fieldStrings := args[1:]

		fields := []generator.InputField{}

		for _, fieldString := range fieldStrings {
			field, err := generator.ParseField(service, name, fieldString)
			if err != nil {
				panic(fmt.Errorf("failed parsing field %s: %w", fieldString, err))
			}

			fields = append(fields, field)
		}

		if parent != "" {
			parentName := generator.TemplateName(parent)
			fields = append(fields, generator.InputField{
				Service:  generator.TemplateName(service),
				Resource: generator.TemplateName(name),
				Name:     generator.TemplateName(parentName.UnderscoreSingular() + "_id"),
				Type:     generator.FieldTypeReferences,
				Required: true,
				Table:    parentName.UnderscorePlural(),
				NotNull:  true,
			})
		}

		input := generator.Input{
			WorkspaceFolder: workspaceFolder,
			HasSearch:       searchField != "",
			Service:         generator.TemplateName(service),
			Resource:        generator.TemplateName(name),
			Fields:          fields,
			SearchField:     searchField,
		}

		if parent != "" {
			v := generator.TemplateName(parent)
			input.Parent = &v
		}

		if err := gen.Generate(cmd.Context(), input); err != nil {
			panic(err)
		}
	},
}

//nolint:gochecknoinits
func init() {
	resourceCmd.Flags().BoolVar(&skipGitCheck, "skip-git", false, "Skip git check for uncommitted changes")
	resourceCmd.Flags().StringVar(&workspaceFolder, "path", ".", "Path to workspace")
	resourceCmd.Flags().StringVar(&searchField, "query-field", "", "Field to search by")
	resourceCmd.Flags().StringVar(&service, "service", "", "Service that the resource belongs to")
	resourceCmd.Flags().StringVar(&parent, "parent", "", "Parent resource")
	resourceCmd.MarkFlagRequired("service") //nolint:errcheck,gosec

	rootCmd.AddCommand(resourceCmd)
}

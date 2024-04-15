package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sparkymat/oxgen/internal/generator"
	"github.com/sparkymat/oxgen/internal/git"
	"github.com/spf13/cobra"
)

var ErrUncommittedChanges = errors.New("uncommitted changes")

var skipGitCheck bool //nolint:gochecknoglobals

var service string //nolint:checknoglobals

var workspaceFolder string //nolint:gochecknoglobals

var searchField string //nolint:gochecknoglobals

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
		fields := args[1:]

		input := generator.GenerateInput{
			WorkspaceFolder: workspaceFolder,
			Name:            name,
			FieldStrings:    fields,
			Service:         service,
			SearchField:     searchField,
		}

		if err := gen.Generate(cmd.Context(), input); err != nil {
			panic(err)
		}
	},
}

//nolint:gochecknoinits
func init() {
	resourceCmd.Flags().BoolVarP(&skipGitCheck, "skip-git", "g", false, "Skip git check for uncommitted changes")
	resourceCmd.Flags().StringVarP(&workspaceFolder, "path", "p", ".", "Path to workspace")
	resourceCmd.Flags().StringVarP(&searchField, "query-field", "q", "", "Field to search by")
	resourceCmd.Flags().StringVarP(&service, "service", "s", "", "Service that the resource belongs to")
	resourceCmd.MarkFlagRequired("service")

	rootCmd.AddCommand(resourceCmd)
}

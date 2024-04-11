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

var name string

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "resource generates a new resource for the project",
	Long:  `resource generates a new resource for the project. `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("Generating resource")

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

		gen := generator.New()

		if err = gen.Generate(cmd.Context(), name); err != nil {
			panic(err)
		}
	},
}

func init() {
	resourceCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the resource to generate")
	resourceCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(resourceCmd)
}

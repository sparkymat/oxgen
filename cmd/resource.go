package cmd

import (
	"errors"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sparkymat/oxgen/internal/git"
	"github.com/spf13/cobra"
)

var ErrUncommittedChanges = errors.New("uncommitted changes")

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
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
}

// Package cmd provides the command-line interface for oxgen.
package cmd

import (
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	Use:   "oxgen",
	Short: "oxgen is a Go web-app project file generator",
	Long: `oxgen generates files for adding new resources to a Go web-app project

Note: oxgen assumes that the project follows certain conventions`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

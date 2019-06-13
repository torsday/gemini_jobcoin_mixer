package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initializeDb)
}

var initializeDb = &cobra.Command{
	Use:   "initDb",
	Short: "initialize db.",
	Long:  `sets up the local relational db`,
	Run: func(cmd *cobra.Command, args []string) {

		Dependencies.PathwayTbl.BuildDbIfNotExists()

	},
}

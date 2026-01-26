package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version = "none"
	Commit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of grocer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version: ", Version, " commit: ", Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

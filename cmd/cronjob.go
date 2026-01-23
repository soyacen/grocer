package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cronjobCmd represents the cronjob command
var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cronjob called")
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
	cronjobCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

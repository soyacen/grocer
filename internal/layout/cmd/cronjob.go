package cmd

import (
	"github.com/soyacen/grocer/internal/layout/internal/cronjob"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		app := fx.New(
			fx.Provide(cronjob.Module),
		)
		return app.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}

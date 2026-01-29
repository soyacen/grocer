package cmd

import (
	"github.com/soyacen/grocer/grocer"
	"github.com/soyacen/grocer/internal/layout/config"
	"github.com/soyacen/grocer/internal/layout/internal/cronjob/cronjob"
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
			ContextModule(ctx),
			config.Module,
			cronjob.Module,
			grocer.Module,
			fx.Invoke(
				func(lc fx.Lifecycle, s *cronjob.Service) {
					lc.Append(fx.StartHook(s.Run))
				},
			),
		)
		return app.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
}

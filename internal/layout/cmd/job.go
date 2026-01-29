package cmd

import (
	"github.com/soyacen/grocer/grocer"
	"github.com/soyacen/grocer/internal/layout/config"
	"github.com/soyacen/grocer/internal/layout/internal/job/job"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var jobCmd = &cobra.Command{
	Use: "job",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		app := fx.New(
			ContextModule(ctx),
			config.Module,
			job.Module,
			grocer.Module,
			fx.Invoke(
				func(lc fx.Lifecycle, s *job.Service) {
					lc.Append(fx.StartHook(s.Run))
				},
			),
		)
		return app.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
}

package cmd

import (
	"context"
	"log/slog"
	"os"

	"go.uber.org/fx"
)

func ContextModule(ctx context.Context) fx.Option {
	return fx.Module("context",
		fx.Provide(func() context.Context { return ctx }),
	)
}

func LoggerModel() fx.Option {
	return fx.Module("logger",
		fx.Provide(func() *slog.Logger {
			return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
		}),
	)
}

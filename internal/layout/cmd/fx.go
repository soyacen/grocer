package cmd

import (
	"context"

	"go.uber.org/fx"
)

func ContextModule(ctx context.Context) fx.Option {
	return fx.Module("context",
		fx.Provide(func() context.Context { return ctx }),
	)
}

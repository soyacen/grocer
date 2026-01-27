package s3x

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"s3x",
	fx.Provide(NewClients),
)

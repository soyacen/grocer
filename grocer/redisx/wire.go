package redisx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongox",
	fx.Provide(NewClients),
)

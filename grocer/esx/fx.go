package esx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"esx",
	fx.Provide(NewClients),
)

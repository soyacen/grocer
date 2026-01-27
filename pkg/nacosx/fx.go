package nacosx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mongox",
	fx.Provide(NewConfigClients, NewNamingClients),
)

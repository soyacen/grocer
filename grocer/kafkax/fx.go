package kafkax

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"kafkax",
	fx.Provide(NewReceivers, NewSenders),
)


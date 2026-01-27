package dbx

import "go.uber.org/fx"

var Module = fx.Module(
	"dbx",
	fx.Provide(NewDBs, NewSqlxDBs),
)

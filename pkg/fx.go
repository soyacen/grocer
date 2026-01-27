package pkg

import (
	"github.com/soyacen/grocer/pkg/dbx"
	"github.com/soyacen/grocer/pkg/esx"
	"github.com/soyacen/grocer/pkg/kafkax"
	"github.com/soyacen/grocer/pkg/mongox"
	"github.com/soyacen/grocer/pkg/nacosx"
	"github.com/soyacen/grocer/pkg/redisx"
	"github.com/soyacen/grocer/pkg/s3x"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"pkg",
	dbx.Module,
	esx.Module,
	kafkax.Module,
	mongox.Module,
	nacosx.Module,
	redisx.Module,
	s3x.Module,
)

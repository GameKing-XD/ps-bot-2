package saver

import "go.uber.org/fx"

var Module = fx.Module("saver", fx.Provide(NewSaver))

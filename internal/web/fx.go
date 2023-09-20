package web

import (
	"github.com/tvanriel/cloudsdk/http"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"web",
	fx.Provide(
		http.AsRouteGroup(NewWeb),
	),
)

package player

import "go.uber.org/fx"

var Module = fx.Module("player",
	fx.Provide(
		NewPlayer,
	),
)

package discord

import (
	"go.uber.org/fx"
	"go.uber.org/multierr"
)

var Module = fx.Module("bot",
	fx.Provide(
		NewDiscord,
	),
	fx.Invoke(
		func(d *DiscordBot) error {
			return multierr.Combine(
				d.Connect(),
				d.AddHandlers(),
			)
		},
	),
)

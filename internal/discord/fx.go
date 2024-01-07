package discord

import (
	"go.uber.org/fx"
)

var Module = fx.Module("bot",
	fx.Provide(
		NewDiscord,
	),
	fx.Invoke(
		AddHandlers,
		ListenQueuedMessages,
	),
)

func AddHandlers(d *DiscordBot) error {
	return d.AddHandlers()
}
func ListenQueuedMessages(d *DiscordBot) error {
	return d.ListenQueuedMessages()
}


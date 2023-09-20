package queues

import "go.uber.org/fx"

var Module = fx.Module("queues", 
        fx.Provide(NewSoundsQueue),
        fx.Provide(NewMessageQueue),
)

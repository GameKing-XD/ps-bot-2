package soundstore

import "go.uber.org/fx"

var Module = fx.Module("soundstore", fx.Provide(NewSoundStore))

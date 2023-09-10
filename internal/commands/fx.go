package commands

import "go.uber.org/fx"

const GROUP_COMMANDS = `group:"commands"`

var Module = fx.Module("commands",
	fx.Provide(
		AsCommand(NewSetPrefixCommand),
		AsCommand(NewListCommand),
		AsCommand(NewJoinCommand),
		AsCommand(NewPSCommand),
		AsCommand(NewSaveCommand),

		fx.Annotate(
			NewCommandExecutor,
			fx.ParamTags(GROUP_COMMANDS),
		),
	),
)

func AsCommand(in any) any {
	return fx.Annotate(
		in,
		fx.As(new(Command)),
		fx.ResultTags(GROUP_COMMANDS),
	)
}

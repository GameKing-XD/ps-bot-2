package commands

import "github.com/tvanriel/ps-bot-2/internal/repositories"

type SetPrefixCommand struct {
	repo *repositories.GuildRepository
}

func NewSetPrefixCommand(g *repositories.GuildRepository) *SetPrefixCommand {

	return &SetPrefixCommand{
		repo: g,
	}
}

func (s *SetPrefixCommand) Name() string {
	return "setprefix"
}

func (s *SetPrefixCommand) Apply(ctx *Context) error {
	if len(ctx.Args) < 1 {
		return nil
	}
	s.repo.UpdatePrefix(ctx.Message.GuildID, ctx.Args[0])
	ctx.Reply("Guild prefix has been updated to " + ctx.Args[0])
	return nil
}

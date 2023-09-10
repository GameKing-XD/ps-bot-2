package commands

import (
	"github.com/tvanriel/ps-bot-2/internal/player"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
)

type PSCommand struct {
	repo   *repositories.GuildRepository
	player *player.Player
}

func NewPSCommand(repo *repositories.GuildRepository, p *player.Player) *PSCommand {
	return &PSCommand{
		player: p,
		repo:   repo,
	}
}

func (ps *PSCommand) Name() string {
	return "ps"
}

func (ps *PSCommand) Apply(ctx *Context) error {
	if len(ctx.Args) < 1 {
		ctx.Reply("Usage: ps <sound>")
		return nil
	}

	ps.player.Queue(ctx.Message.GuildID, ctx.Args[0])

	return nil
}

package commands

import (
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
)

type PSCommand struct {
	repo   *repositories.GuildRepository
        queue *queues.SoundsQueue
}

func NewPSCommand(repo *repositories.GuildRepository, queue *queues.SoundsQueue) *PSCommand {
	return &PSCommand{
                queue: queue,
		repo:   repo,
	}
}

func (ps *PSCommand) Name() string {
	return "ps"
}

func (p *PSCommand) SkipsPrefix() bool {
	return false
}
func (ps *PSCommand) Apply(ctx *Context) error {
	if len(ctx.Args) < 1 {
		ctx.Reply("Usage: ps <sound>")
		return nil
	}

        err := ps.queue.Append(ctx.Message.GuildID, ctx.Args[0])
        if err != nil {
                return err
        }

	return nil
}

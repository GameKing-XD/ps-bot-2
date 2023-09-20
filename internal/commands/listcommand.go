package commands

import "github.com/tvanriel/ps-bot-2/internal/soundstore"

type ListCommand struct {
	store *soundstore.SoundStore
}

func NewListCommand(store *soundstore.SoundStore) *ListCommand {
	return &ListCommand{
		store: store,
	}
}
func (l *ListCommand) Name() string {
	return "pslist"
}

func (l *ListCommand) SkipsPrefix() bool {
	return false
}
func (l *ListCommand) Apply(ctx *Context) error {

	ctx.ReplyList(l.store.List(ctx.Message.GuildID))

	return nil
}

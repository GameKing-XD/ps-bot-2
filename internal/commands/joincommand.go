package commands

import "github.com/tvanriel/ps-bot-2/internal/repositories"

type JoinCommand struct {
	repo *repositories.GuildRepository
}

func NewJoinCommand(repo *repositories.GuildRepository) *JoinCommand {
	return &JoinCommand{
		repo: repo,
	}
}

func (j *JoinCommand) Name() string {
	return "join"
}

func (j *JoinCommand) SkipsPrefix() bool {
	return false
}

func (j *JoinCommand) Apply(ctx *Context) error {
	state, err := ctx.Session.State.VoiceState(ctx.Message.GuildID, ctx.Message.Author.ID)
	if err != nil {
		return err
	}
	if state.ChannelID == "" {
		ctx.Reply("You must be in a voicechannel to do this.")
		return nil
	}

	_, err = ctx.Session.ChannelVoiceJoin(ctx.Message.GuildID, state.ChannelID, false, true)

	if err != nil {
		ctx.Reply(err.Error())
		return err
	}

	j.repo.JoinVoiceChannel(ctx.Message.GuildID, state.ChannelID)

	return nil
}

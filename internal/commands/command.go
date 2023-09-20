package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/multierr"
)

type Command interface {
	Name() string
	Apply(*Context) error
	SkipsPrefix() bool
}

type Context struct {
	Message *discordgo.Message
	Content string
	Args    []string
	Session *discordgo.Session
}

func (ctx *Context) Reply(s string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbedReply(
		ctx.Message.ChannelID,

		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Result",
					Value: s,
				},
			},
		},
		ctx.Reference(),
	)

}

func (ctx *Context) Error(err error) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendEmbedReply(
		ctx.Message.ChannelID,
		&discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: err.Error(),
				},
			},
			Color: 0xFF0000,
		},
		ctx.Reference(),
	)
}

func (ctx *Context) Reference() *discordgo.MessageReference {

	return &discordgo.MessageReference{
		MessageID: ctx.Message.ID,
		ChannelID: ctx.Message.ChannelID,
		GuildID:   ctx.Message.GuildID,
	}
}

const DISCORD_MSG_MAX_LEN = 2000

func (ctx *Context) ReplyList(s []string) ([]*discordgo.Message, error) {
	if len(s) == 0 {
		return []*discordgo.Message{}, nil

	}

	itemTpl := "`%s`\n"
	templated := make([]string, len(s))
	for i := range s {
		templated = append(templated, fmt.Sprintf(itemTpl, s[i]))
	}
	var sb strings.Builder
	totalLength := 0

	var sentMessages []*discordgo.Message
	var err error
	for i := range templated {
		if totalLength+len(templated[i]) > DISCORD_MSG_MAX_LEN {
			msg, msgerr := ctx.Session.ChannelMessageSendReply(
				ctx.Message.ChannelID,
				sb.String(),
				ctx.Reference(),
			)
			sentMessages = append(sentMessages, msg)
			err = multierr.Append(err, msgerr)
			sb.Reset()
			totalLength = 0
		}

		totalLength += len(templated[i])
		sb.WriteString(templated[i])
	}
	if totalLength > 0 {
		msg, msgerr := ctx.Session.ChannelMessageSendReply(
			ctx.Message.ChannelID,
			sb.String(),
			ctx.Reference(),
		)
		sentMessages = append(sentMessages, msg)
		err = multierr.Append(err, msgerr)
	}
	return sentMessages, err

}

package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"strings"
)

type Executor struct {
	commands []Command
	log      *zap.Logger
}

func NewCommandExecutor(commands []Command, log *zap.Logger) *Executor {
	return &Executor{
		commands: commands,
		log:      log,
	}
}
func (e *Executor) HasMatch(trigger string, message string) bool {
	for i := range e.commands {
		if !e.commands[i].SkipsPrefix() && HasCommandPrefix(trigger, e.commands[i].Name(), message) {
			return true
		}
		if e.commands[i].SkipsPrefix() && e.commands[i].Name() == message {
			return true
		}
	}
	return false
}
func (e *Executor) Apply(trigger string, message *discordgo.Message, s *discordgo.Session) {

	for i := range e.commands {
		skipsPrefix := e.commands[i].SkipsPrefix()
		messageMatchesName := (e.commands[i].Name() == message.Content)
		commandPrefix := (!skipsPrefix) && HasCommandPrefix(trigger, e.commands[i].Name(), message.Content)

		if skipsPrefix && messageMatchesName {
			go func(cmd Command) {
				ctx := &Context{
					Message: message,
					Args:    []string{message.Content},
					Session: s,
					Content: message.Content,
				}

				err := cmd.Apply(ctx)
				if err != nil {
					e.log.Error(
						"Command failed",
						zap.NamedError("err", err),
						zap.String("username", message.Author.Username),
						zap.String("guild", message.GuildID),
						zap.String("channel", message.ChannelID),
						zap.String("message", message.ID),
						zap.String("content", message.Content),
						zap.String("url", messagePermaUrl(message.GuildID, message.ChannelID, message.ID)),
					)

					_, err1 := ctx.Error(err)
					if err1 != nil {
						e.log.Error(
							"Failed to report command reply error to discord",
							zap.NamedError("orig", err),
							zap.NamedError("err", err1),
							zap.String("username", message.Author.Username),
							zap.String("guild", message.GuildID),
							zap.String("channel", message.ChannelID),
							zap.String("message", message.ID),
							zap.String("content", message.Content),
							zap.String("url", messagePermaUrl(message.GuildID, message.ChannelID, message.ID)),
						)
					}
				}

			}(e.commands[i])
			continue
		}

		if commandPrefix {
			content := StripPrefix(trigger, e.commands[i].Name())(message.Content)
			args := SplitArgs(content)
			go func(cmd Command) {
				ctx := &Context{
					Message: message,
					Content: content,
					Args:    args,
					Session: s,
				}

				err := cmd.Apply(ctx)

				if err != nil {
					e.log.Error(
						"Command failed",
						zap.NamedError("err", err),
						zap.String("username", message.Author.Username),
						zap.String("guild", message.GuildID),
						zap.String("channel", message.ChannelID),
						zap.String("message", message.ID),
						zap.String("content", message.Content),
						zap.String("url", messagePermaUrl(message.GuildID, message.ChannelID, message.ID)),
					)

					_, err1 := ctx.Error(err)
					if err1 != nil {
						e.log.Error(
							"Failed to report command reply error to discord",
							zap.NamedError("orig", err),
							zap.NamedError("err", err1),
							zap.String("username", message.Author.Username),
							zap.String("guild", message.GuildID),
							zap.String("channel", message.ChannelID),
							zap.String("message", message.ID),
							zap.String("content", message.Content),
							zap.String("url", messagePermaUrl(message.GuildID, message.ChannelID, message.ID)),
						)
					}
				}
			}(e.commands[i])
		}

	}

}

func messagePermaUrl(guild string, channel string, id string) string {
	return strings.Join(
		[]string{
			"https://discord.com/channels/",
			guild,
			"/",
			channel,
			"/",
			id,
		},
		"",
	)
}

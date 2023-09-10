package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Executor struct {
	commands []Command
}

func NewCommandExecutor(commands []Command) *Executor {
	return &Executor{
		commands: commands,
	}
}
func (e *Executor) HasMatch(trigger string, message string) bool {
	for i := range e.commands {
		if HasCommandPrefix(trigger, e.commands[i].Name())(message) {
			return true
		}
	}
	return false
}
func (e *Executor) Apply(trigger string, message *discordgo.Message, s *discordgo.Session) {

	for i := range e.commands {
		if HasCommandPrefix(trigger, e.commands[i].Name())(message.Content) {
			content := StripPrefix(trigger, e.commands[i].Name())(message.Content)
			args := SplitArgs(content)
			go e.commands[i].Apply(
				&Context{
					Message: message,
					Content: content,
					Args:    args,
					Session: s,
				},
			)
		}
	}

}

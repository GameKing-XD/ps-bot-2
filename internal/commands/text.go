package commands

import "strings"

// Strip the prefix from a command
func StripPrefix(trigger string, command string) func(string) string {
	prefixLen := len(trigger) + len(command)
	return func(message string) string {
		return strings.TrimSpace(message[prefixLen:])
	}
}

func HasCommandPrefix(trigger string, command string, message string) bool {
	if message == trigger+command {
		return true
	}
	return strings.HasPrefix(message, trigger+command+" ")
}

func SplitArgs(s string) []string {
	return strings.Split(s, " ")
}

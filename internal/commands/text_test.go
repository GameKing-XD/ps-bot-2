package commands_test

import (
	"testing"

	"github.com/tvanriel/ps-bot-2/internal/commands"
	assert "gotest.tools/v3/assert"
)

func TestStripPrefix(t *testing.T) {
	assert.Equal(t, commands.StripPrefix("!", "ps")("!ps test"), "test")
	assert.Equal(t, commands.StripPrefix("!", "ps")("!ps       test"), "test")
	assert.Equal(t, commands.StripPrefix("!", "ps")("!pstest"), "est")
}

func TestHasPrefix(t *testing.T) {
	assert.Equal(t, commands.HasCommandPrefix("!", "ps")("!ps test"), true)
	assert.Equal(t, commands.HasCommandPrefix("!", "ps")("!ps      test"), true)
	assert.Equal(t, commands.HasCommandPrefix("!", "ps")("!pstest"), false)

	assert.Equal(t, commands.HasCommandPrefix("!", "ps")("!pslist"), false)
}

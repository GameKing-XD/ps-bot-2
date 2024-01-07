package commands

import (
	"encoding/json"
	"strings"

	"github.com/tvanriel/ps-bot-2/internal/metrics"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/saver"
	"go.uber.org/zap"
)

var kerstterms = []string{
	"kerst",
	"christmas",
	"christ",
}

func kerstbtfo(banned []string) func(s string) bool {

	return func(s string) bool {
		for _, word := range banned {
			if s == word {
				return true
			}
		}
		return false
	}
}

type SaveCommand struct {
	secretName string
	bucketName string
	log        *zap.Logger
	saver      *saver.Saver
        metrics    *metrics.MetricsCollector
}

func (s *SaveCommand) SkipsPrefix() bool {
	return false
}

func NewSaveCommand(l *zap.Logger, saver *saver.Saver, metrics *metrics.MetricsCollector) *SaveCommand {
	return &SaveCommand{
		log:   l,
		saver: saver,
                metrics: metrics,
	}
}

func (s *SaveCommand) Name() string {
	return "save"
}

func (s *SaveCommand) Apply(ctx *Context) error {

	if len(ctx.Args) < 1 {
		ctx.Reply("Usage: save <name> - Saves the attachment as a ps command")
		return nil
	}

	if len(ctx.Message.Attachments) != 1 {
		ctx.Reply("You must provide an attachment")
		return nil
	}

	url := ctx.Message.Attachments[0].URL
	guildId := ctx.Message.GuildID
	soundName := ctx.Args[0]

	s.log.Info("Pushing job to Kubernetes",
		zap.String("url", url),
		zap.String("guildId", guildId),
		zap.String("soundName", soundName),
	)

	amqpBody, err := json.Marshal(queues.QueuedMessage{
		ChannelID: ctx.Message.ChannelID,
		Content: strings.Join([]string{
			"Saved sound ",
			soundName,
			".",
		}, ""),
		GuildID: guildId,
	})
	if err != nil {
		return err
	}
        s.metrics.RegisterPlaySound(guildId, soundName)

	return s.saver.Save(saver.SaveParams{
		ChannelID:   ctx.Message.ChannelID,
		GuildID:     guildId,
		SoundName:   soundName,
		URL:         url,
		TextMessage: string(amqpBody),
	})
}

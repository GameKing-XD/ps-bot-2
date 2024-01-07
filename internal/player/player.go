package player

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/tvanriel/ps-bot-2/internal/metrics"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ErrNoSuchConnection = errors.New("no connection to guild voice channel")

type Player struct {
	queue *queues.SoundsQueue
	log   *zap.Logger
	repo  *repositories.GuildRepository
	store *soundstore.SoundStore
        Metrics *metrics.MetricsCollector
}

type NewPlayerParams struct {
        fx.In
Log *zap.Logger
        Repo *repositories.GuildRepository
        Store *soundstore.SoundStore
        Queue *queues.SoundsQueue
        Metrics *metrics.MetricsCollector
}

func NewPlayer(params NewPlayerParams) *Player {
	return &Player{
		queue: params.Queue,
		log:   params.Log.Named("player"),
		repo:  params.Repo,
		store: params.Store,
                Metrics: params.Metrics,
	}
}

func (p *Player) Connect(ses *discordgo.Session, guildId string) error {

	log := p.log.With(zap.String("guildId", guildId))
	conn, err := p.attemptConnect(ses, guildId)
	if err != nil {
		return err
	}
	log.Info("Connected to voicechannel")

	msgs := p.queue.Consume(guildId)

	go func() {
		log.Info("Looping for messages")
		for m := range msgs {
			s := m.Sound
			log.Info("playing sound",
				zap.String("sound", s),
			)

			reader, err := p.store.Find(guildId, s)
			if err != nil {
				log.Warn("Could not play sound, store does not contain sound",
					zap.String("sound", s),
				)
				continue
			}
			buf, err := loadSound(reader)
			if err != nil {

				continue
			}

                        p.Metrics.PlaySound(guildId, s)
			err = conn.Speaking(true)
			if err != nil {
				continue
			}
			for i := range buf {
				conn.OpusSend <- buf[i]
			}
			_ = conn.Speaking(false)
		}
	}()

	return nil
}

func (p *Player) attemptConnect(ses *discordgo.Session, guildId string) (*discordgo.VoiceConnection, error) {

	conn, ok := ses.VoiceConnections[guildId]
	if ok {
		return conn, nil
	}

	channelId := p.repo.GetVoiceChannel(guildId)

	if channelId == "" {
		return nil, ErrNoSuchConnection
	}
	return ses.ChannelVoiceJoin(guildId, channelId, false, true)

}

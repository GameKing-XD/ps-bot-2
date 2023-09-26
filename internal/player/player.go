package player

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/zap"
)

var ErrNoSuchConnection = errors.New("no connection to guild voice channel")

type Player struct {
	queue *queues.SoundsQueue
	log   *zap.Logger
	repo  *repositories.GuildRepository
	store *soundstore.SoundStore
}

func NewPlayer(log *zap.Logger, repo *repositories.GuildRepository, store *soundstore.SoundStore, queue *queues.SoundsQueue) *Player {
	return &Player{
		queue: queue,
		log:   log,
		repo:  repo,
		store: store,
	}
}

func (p *Player) Connect(ses *discordgo.Session, guildId string) error {

	conn, err := p.attemptConnect(ses, guildId)
	if err != nil {
		return err
	}
	p.log.Info("Connected to voicechannel", zap.String("guild", guildId))

	msgs, err := p.queue.Consume(guildId)

	if err != nil {
		return err
	}
	go func() {

		for m := range msgs {
			s := m.Sound
			p.log.Info("playing sound",
				zap.String("sound", s),
				zap.String("guildId", guildId),
			)

			reader, err := p.store.Find(guildId, s)
			if err != nil {
				p.log.Warn("Could not play sound, store does not contain sound",
					zap.String("sound", s),
					zap.String("guildId", guildId),
				)
				continue
			}
			buf, err := loadSound(reader)
			if err != nil {

				continue
			}

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

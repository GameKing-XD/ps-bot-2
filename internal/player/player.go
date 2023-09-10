package player

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/zap"
)

var ErrNoSuchConnection = errors.New("no connection to guild voice channel")

type Player struct {
	Queues map[string](chan string)
	log    *zap.Logger
	repo   *repositories.GuildRepository
	store  *soundstore.SoundStore
}

func NewPlayer(log *zap.Logger, repo *repositories.GuildRepository, store *soundstore.SoundStore) *Player {
	return &Player{
		Queues: make(map[string]chan string),
		log:    log,
		repo:   repo,
		store:  store,
	}
}

func (p *Player) Connect(ses *discordgo.Session, guildId string) error {

	p.Queues[guildId] = make(chan string, 50)

	conn, err := p.attemptConnect(ses, guildId)
	if err != nil {
		return err
	}

	go func() {
		for s := range p.Queues[guildId] {
			reader, err := p.store.Find(guildId, s)
			if err != nil {
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

func (p *Player) Queue(guildId string, soundName string) {
	p.Queues[guildId] <- soundName
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

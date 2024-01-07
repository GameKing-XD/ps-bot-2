package queues

import (
	"context"
	"encoding/json"

	"github.com/tvanriel/cloudsdk/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SoundsQueue struct {
	Log   *zap.Logger
	Redis *redis.RedisClient
}

const REDIS_SOUND_CHAN = "sounds"

type NewSoundsQueueParams struct {
	fx.In

	Redis *redis.RedisClient
	Log   *zap.Logger
}

func NewSoundsQueue(params NewMessageQueueParams) *SoundsQueue {
	s := &SoundsQueue{
		Log:   params.Log.Named("sound-queue"),
		Redis: params.Redis,
	}
	return s

}

type playSoundMessage struct {
	GuildID string `json:"guild"`
	Sound   string `json:"sound"`
}

func (sq *SoundsQueue) Consume(guildId string) chan playSoundMessage {

	log := sq.Log.With(zap.String("guildId", guildId))
	pubsub := sq.Redis.Conn().Subscribe(context.Background(), REDIS_SOUND_CHAN+"-"+guildId)
	rch := pubsub.Channel()
	mch := make(chan playSoundMessage)
	go func() {

		for rm := range rch {
			log.Debug("Received message from Redis")

			var qm playSoundMessage
			err := json.Unmarshal([]byte(rm.Payload), &qm)
			if err != nil {
				log.Error("Cannot Unmarshal sounds from redis subscribe", zap.Error(err))
				continue
			}
			log.Debug("Acknowledged message", zap.String("sound", qm.Sound))

			mch <- qm

		}
	}()
	return mch
}
func (sq *SoundsQueue) Append(guild, sound string) error {

	b, err := json.Marshal(playSoundMessage{
		GuildID: guild,
		Sound:   sound,
	})

	if err != nil {
		return err
	}

	cmd := sq.Redis.Conn().Publish(context.Background(), REDIS_SOUND_CHAN+"-"+guild, b)
	return cmd.Err()

}

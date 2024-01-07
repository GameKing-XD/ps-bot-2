package queues

import (
	"context"
	"encoding/json"

	"github.com/tvanriel/cloudsdk/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MessageQueue struct {
	Log   *zap.Logger
	Redis *redis.RedisClient
}

type NewMessageQueueParams struct {
	fx.In

	Redis *redis.RedisClient
	Log   *zap.Logger
}

const REDIS_CHAN = "messages"

func NewMessageQueue(p NewMessageQueueParams) *MessageQueue {
	return &MessageQueue{
		Log:   p.Log.Named("messagequeue"),
		Redis: p.Redis,
	}

}

type QueuedMessage struct {
	ChannelID string `json:"channelId"`
	GuildID   string `json:"guildId"`
	Content   string `json:"content"`
}

func (m *MessageQueue) Consume() chan *QueuedMessage {
	pubsub := m.Redis.Conn().Subscribe(context.Background(), REDIS_CHAN)
	rch := pubsub.Channel()
	mch := make(chan *QueuedMessage)
        go func() {

	for rm := range rch {
		var qm *QueuedMessage
		err := json.Unmarshal([]byte(rm.Payload), qm)
		if err != nil {
			m.Log.Error("Cannot Unmarshal message from redis subscribe", zap.Error(err))
			continue
		}

		mch <- qm

	}
        }()
	return mch
}

func (m *MessageQueue) Append(msg QueuedMessage) error {
	b, err := json.Marshal(msg)

	if err != nil {
		return err
	}
	cmd := m.Redis.Conn().Publish(context.Background(), REDIS_CHAN, b)

	return cmd.Err()
}

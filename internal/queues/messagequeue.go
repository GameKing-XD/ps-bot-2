package queues

import (
	"encoding/json"

	"github.com/tvanriel/cloudsdk/amqp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MessageQueue struct {
	Amqp *amqp.Connection
	Log  *zap.Logger
}

type NewMessageQueueParams struct {
	fx.In

	Amqp *amqp.Connection
	Log  *zap.Logger
}

func NewMessageQueue(p NewMessageQueueParams) *MessageQueue {
	return &MessageQueue{
		Amqp: p.Amqp,
		Log:  p.Log.Named("messagequeue"),
	}

}

type QueuedMessage struct {
	ChannelID string `json:"channelId"`
	GuildID   string `json:"guildId"`
	Content   string `json:"content"`
}

func (m *MessageQueue) Consume() (chan QueuedMessage, error) {

	out := make(chan QueuedMessage)
	ch, err := m.Amqp.Channel()
	if err != nil {
		return nil, err
	}

	ch.QueueDeclare("messages", false, false, false, true, nil)
	if err != nil {
		return nil, err
	}
	msgs, err := ch.Consume("messages", "", true, false, false, true, nil)

	if err != nil {
		return out, err
	}
	go func() {
		for x := range msgs {
			m.Log.Info("Recv on AMQP Messages queue", zap.String("body", string(x.Body)))
			message := QueuedMessage{}
			err := json.Unmarshal(x.Body, &message)
			if err != nil {
				m.Log.Error("cannot unmarshal message from messagequeue", zap.Error(err))
				continue
			}
			out <- message

		}
	}()

	return out, nil
}

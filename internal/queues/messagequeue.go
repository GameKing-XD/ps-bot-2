package queues

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type MessageQueue struct {
	Ch  *amqp091.Channel
	log *zap.Logger
}

func NewMessageQueue(amqp *amqp091.Connection, l *zap.Logger) (*MessageQueue, error) {
	ch, err := amqp.Channel()
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		Ch:  ch,
		log: l,
	}, nil

}

type QueuedMessage struct {
	ChannelID string `json:"channelId"`
	GuildID   string `json:"guildId"`
	Content   string `json:"content"`
}

func (m *MessageQueue) Append() {

}

func (m *MessageQueue) Consume() (chan QueuedMessage, error) {

	out := make(chan QueuedMessage)
	_, err := m.Ch.QueueDeclare("messages", false, false, false, true, nil)
	if err != nil {
		return nil, err
	}
	msgs, err := m.Ch.Consume("messages", "", true, false, false, true, nil)

	if err == nil {
		go func() {
			for x := range msgs {
				m.log.Info("Recv on AMQP Messages queue", zap.String("body", string(x.Body)))
				message := QueuedMessage{}
				err := json.Unmarshal(x.Body, &message)
				if err != nil {
					m.log.Error("cannot unmarshal message from messagequeue", zap.Error(err))
					continue
				}
				out <- message

			}
		}()
	}
	return out, err

}

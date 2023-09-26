package queues

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tvanriel/cloudsdk/amqp"
	"go.uber.org/zap"
)

type SoundsQueue struct {
	Amqp *amqp.Connection
	Log  *zap.Logger
}

func NewSoundsQueue(amqp *amqp.Connection, log *zap.Logger) *SoundsQueue {
	s := &SoundsQueue{
		Amqp: amqp,
		Log:  log.Named("sound-queue"),
	}
	return s

}

type playSoundMessage struct {
	GuildID string `json:"guild"`
	Sound   string `json:"sound"`
}

func (s *SoundsQueue) Append(guildId string, sound string) error {
	ch, err := s.Amqp.Channel()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare("sounds-"+guildId, false, false, false, true, nil)
	if err != nil {
		return err
	}

	body := &playSoundMessage{
		GuildID: guildId,
		Sound:   sound,
	}
	marshalled, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(context.Background(), "", "sounds-"+guildId, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        marshalled,
	})
}

func (s *SoundsQueue) Consume(guildId string) (chan playSoundMessage, error) {
	ch, err := s.Amqp.Channel()
	if err != nil {
		return nil, err
	}
	out := make(chan playSoundMessage)
	_, err = ch.QueueDeclare("sounds-"+guildId, false, false, false, true, nil)
	if err != nil {
		return nil, err
	}
	msgs, err := ch.Consume("sounds-"+guildId, "", true, false, false, true, nil)

	if err != nil {
		return nil, err
	}
	go func() {
		for {
			for m := range msgs {
				message := playSoundMessage{}
				err := json.Unmarshal(m.Body, &message)
				if err != nil {
					s.Log.Error("cannot unmarshal message from queue", zap.Error(err))
					continue
				}
				out <- message
			}
		}
	}()
	return out, err
}

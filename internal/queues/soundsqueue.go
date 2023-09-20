package queues

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)


type SoundsQueue struct{
	Ch    *amqp091.Channel
        log   *zap.Logger
        amqp *amqp091.Connection
}

func NewSoundsQueue(amqp *amqp091.Connection, log *zap.Logger) (*SoundsQueue, error) {
        s := &SoundsQueue{
                amqp: amqp,
                Ch: nil,
                log: log,
        }
        err := s.Connect()
        return s, err

}

type playSoundMessage struct{
        GuildID string `json:"guild"`
        Sound   string `json:"sound"`
}

func (s *SoundsQueue) Connect() error {
        if s.Ch != nil && !s.Ch.IsClosed() {
                return nil
        }

        ch, err := s.amqp.Channel()
        if err != nil {
                return err
        }
        s.Ch = ch
        return nil
}

func (s *SoundsQueue) Append(guildId string, sound string) error {
        err := s.Connect()
        if err != nil {
                return err
        }

        _, err = s.Ch.QueueDeclare("sounds-" + guildId, false, false, false, true, nil)
        if err != nil {
                return err
        }

        body := &playSoundMessage{
                GuildID: guildId,
                Sound: sound,
        }
		marshalled, err := json.Marshal(body)
        if err != nil {
                return err
        }

	return s.Ch.PublishWithContext(context.Background(), "", "sounds-" + guildId, false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        marshalled,
		})
}

func (s *SoundsQueue) Consume(guildId string) (chan playSoundMessage, error) {
        err := s.Connect()
        if err != nil {
                return nil, err
        }
        out := make(chan playSoundMessage)
        _, err = s.Ch.QueueDeclare("sounds-" + guildId, false, false, false, true, nil)
        if err != nil {
               return nil, err 
        }
        msgs, err := s.Ch.Consume("sounds-" + guildId, "", true, false, false, true, nil)
        
        if err == nil {
                go func() {
                        for m := range msgs {
                                message := playSoundMessage{}
                                err := json.Unmarshal(m.Body, &message)
                                if err != nil {
                                        s.log.Error("cannot unmarshal message from queue", zap.Error(err))
                                }
                                out <- message
                        }
                }()
        }
        return out, err
}

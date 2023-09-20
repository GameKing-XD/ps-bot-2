package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/tvanriel/ps-bot-2/internal/commands"
	"github.com/tvanriel/ps-bot-2/internal/player"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"go.uber.org/zap"
)

type DiscordBot struct {
	conn   *discordgo.Session
	log    *zap.Logger
	ready  bool
	repo   *repositories.GuildRepository
	player *player.Player
	exe    *commands.Executor
        queue  *queues.MessageQueue
}

func NewDiscord(config *Configuration, log *zap.Logger, repo *repositories.GuildRepository, exe *commands.Executor, p *player.Player, queue *queues.MessageQueue) (*DiscordBot, error) {
	ses, err := discordgo.New("Bot " + config.BotToken)

	if err != nil {
		return nil, err
	}
	ses.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	return &DiscordBot{
		conn:   ses,
		repo:   repo,
		ready:  false,
		log:    log,
		exe:    exe,
		player: p,
                queue:  queue,
	}, nil
}

func (d *DiscordBot) AddHandlers() error {
	d.conn.AddHandler(messagehandler(d))
	d.conn.AddHandler(ready(d))
	d.conn.AddHandler(guildCreate(d))
	return nil
}

func (d *DiscordBot) ListenQueuedMessages() error {
        msgs, err := d.queue.Consume()
        if err != nil {
                return err
        }
        go func ()  {
                
                for m := range msgs {
                        d.log.Info("Send message from AMQP chan", 
                                zap.String("channel", m.ChannelID), 
                                zap.String("content", m.Content),
                        )

                        if m.ChannelID == "" || m.Content == "" {
                                d.log.Error("Invalid message request from AMQP chan",
                                        zap.String("channel", m.ChannelID),
                                        zap.String("content", m.Content),
                                )
                                return
                        }
                        _, err := d.conn.ChannelMessageSend(m.ChannelID, m.Content)
                        if err != nil {
                                d.log.Error("error while listening to queued messages", 
                                        zap.Error(err),
                                        zap.String("channel", m.ChannelID),
                                        zap.String("content", m.Content),
                                )
                        }

                }
        }()
        return nil
}

func (d *DiscordBot) Connect() error {
	return d.conn.Open()
}

func (d *DiscordBot) JoinVoiceChannels() {
	d.repo.GetVoiceChannels()
}

func (d *DiscordBot) PlayVoiceCommand(s *discordgo.Session, sound string, guildId string) {
	fmt.Println(sound, guildId)

}

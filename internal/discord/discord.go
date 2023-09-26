package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tvanriel/ps-bot-2/internal/commands"
	"github.com/tvanriel/ps-bot-2/internal/player"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DiscordBot struct {
	Ready  bool
	Conn   *discordgo.Session
	Log    *zap.Logger
	Repo   *repositories.GuildRepository
	Player *player.Player
	Exe    *commands.Executor
	Queue  *queues.MessageQueue
}

type NewDiscordParams struct {
	fx.In

	Config *Configuration
	Conn   *discordgo.Session
	Log    *zap.Logger
	Repo   *repositories.GuildRepository
	Player *player.Player
	Exe    *commands.Executor
	Queue  *queues.MessageQueue
}

func NewDiscord(p NewDiscordParams) (*DiscordBot, error) {
	ses, err := discordgo.New("Bot " + p.Config.BotToken)

	if err != nil {
		return nil, err
	}
	ses.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	return &DiscordBot{
		Conn:   ses,
		Repo:   p.Repo,
		Ready:  false,
		Log:    p.Log.Named("discord"),
		Exe:    p.Exe,
		Player: p.Player,
		Queue:  p.Queue,
	}, nil
}

func (d *DiscordBot) AddHandlers() error {
	d.Conn.AddHandler(messagehandler(d))
	d.Conn.AddHandler(ready(d))
	d.Conn.AddHandler(guildCreate(d))
	return nil
}

func (d *DiscordBot) ListenQueuedMessages() error {
	msgs, err := d.Queue.Consume()
	if err != nil {
		return err
	}
	go func() {

		for m := range msgs {
			d.Log.Info("Send message from AMQP chan",
				zap.String("channel", m.ChannelID),
				zap.String("content", m.Content),
			)

			if m.ChannelID == "" || m.Content == "" {
				d.Log.Error("Invalid message request from AMQP chan",
					zap.String("channel", m.ChannelID),
					zap.String("content", m.Content),
				)
				return
			}

			content := escapeDiscordMessage(m.Content)

			_, err := d.Conn.ChannelMessageSend(m.ChannelID, content)
			if err != nil {
				d.Log.Error("error while listening to queued messages",
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
	return d.Conn.Open()
}

func (d *DiscordBot) JoinVoiceChannels() {
	d.Repo.GetVoiceChannels()
}

func (d *DiscordBot) PlayVoiceCommand(s *discordgo.Session, sound string, guildId string) {
	fmt.Println(sound, guildId)

}
func escapeDiscordMessage(s string) string {
	s = strings.ReplaceAll(s, "@", "")
	s = strings.ReplaceAll(s, "#", "")

	return s
}

package app

import (
	"github.com/tvanriel/cloudsdk/amqp"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/mysql"
	"github.com/tvanriel/cloudsdk/s3"
	"github.com/tvanriel/ps-bot-2/internal/commands"
	"github.com/tvanriel/ps-bot-2/internal/config"
	"github.com/tvanriel/ps-bot-2/internal/discord"
	"github.com/tvanriel/ps-bot-2/internal/player"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"github.com/tvanriel/ps-bot-2/internal/web"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func DiscordBot() {
	fx.New(
		fx.Provide(
			config.ViperConfiguration,
			config.MySQLConfiguration,
			config.DiscordConfiguration,
			config.LoggingConfiguration,
			config.StorageConfiguration,
			config.S3Configuration,
			config.KubernetesConfiguration,
			config.SaverConfiguration,
			config.AmqpConfiguration,
		),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		mysql.Module,
		logging.Module,
		player.Module,
		discord.Module,
		repositories.Module,
		commands.Module,
		soundstore.Module,
		s3.Module,
		kubernetes.Module,
		amqp.Module,
		queues.Module,
		fx.Invoke(func(_ *discord.DiscordBot) {}),
	).Run()
}

func Web() {
	fx.New(
		fx.Provide(
			config.ViperConfiguration,
			config.MySQLConfiguration,
			config.HttpConfiguration,
			config.AmqpConfiguration,
			config.LoggingConfiguration,
			config.S3Configuration,
			config.StorageConfiguration,
		),
		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l}
		}),
		mysql.Module,
		http.Module,
		web.Module,
		amqp.Module,
		logging.Module,
		repositories.Module,
		soundstore.Module,
		queues.Module,
		s3.Module,
		fx.Invoke(http.Listen),
	).Run()
}

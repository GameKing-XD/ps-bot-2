package app

import (
	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/mysql"
	"github.com/tvanriel/cloudsdk/s3"
	"github.com/tvanriel/ps-bot-2/internal/commands"
	"github.com/tvanriel/ps-bot-2/internal/config"
	"github.com/tvanriel/ps-bot-2/internal/discord"
	"github.com/tvanriel/ps-bot-2/internal/player"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/fx"
)

func Run() {
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
		),
		mysql.Module,
		logging.Module,
		player.Module,
		discord.Module,
		repositories.Module,
		commands.Module,
		soundstore.Module,
		s3.Module,
		kubernetes.Module,
	).Run()
}

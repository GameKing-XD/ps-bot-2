package saver

import (
	"strings"

	"github.com/tvanriel/cloudsdk/kubernetes"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Configuration struct {
	BucketName string
	SecretName string
}

type NewSaverParams struct {
	fx.In

	Kubernetes    *kubernetes.KubernetesClient
	Configuration *Configuration
	Logger        *zap.Logger
	//Metrics *metrics.MetricsCollector
}

type Saver struct {
	Kubernetes    *kubernetes.KubernetesClient
	Configuration *Configuration
	Logger        *zap.Logger
	//Metrics *metrics.MetricsCollector
}

func NewSaver(s NewSaverParams) *Saver {
	return &Saver{
		Logger:        s.Logger.Named("saver"),
		Kubernetes:    s.Kubernetes,
		Configuration: s.Configuration,
		//Metrics: s.Metrics,
	}
}

type SaveParams struct {
	ChannelID   string
	GuildID     string
	SoundName   string
	URL         string
	TextMessage string
}

func (s SaveParams) Target(config *Configuration) string {
	return strings.Join([]string{
		config.BucketName, "/", s.GuildID,
	}, "")
}

func (s *Saver) Save(params SaveParams) error {
	return s.Kubernetes.RunJob(convertJob(params, s.Configuration))
}

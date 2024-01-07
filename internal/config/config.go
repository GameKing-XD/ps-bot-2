package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/tvanriel/cloudsdk/http"
	"github.com/tvanriel/cloudsdk/kubernetes"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/mysql"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/cloudsdk/redis"
	"github.com/tvanriel/cloudsdk/s3"
	"github.com/tvanriel/ps-bot-2/internal/discord"
	"github.com/tvanriel/ps-bot-2/internal/saver"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
)

type Configuration struct {
	MySQL      mysql.Configuration      `mapstructure:"mysql"`
	Logging    logging.Configuration    `mapstructure:"log"`
	Discord    discord.Configuration    `mapstructure:"discord"`
	S3         s3.Configuration         `mapstructure:"s3"`
	Storage    soundstore.Configuration `mapstructure:"storage"`
	Kubernetes kubernetes.Configuration `mapstructure:"kubernetes"`
	Saver      saver.Configuration      `mapstructure:"saver"`
	Http       http.Configuration       `mapstructure:"http"`
	Redis      redis.Configuration      `mapstructure:"redis"`
        Prometheus prometheus.Configuration `mapstructure:"prometheus"`
}

func ViperConfiguration() (Configuration, error) {
	var config Configuration
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/discordbot")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Configuration{}, err
	}

	return config, err
}

func MySQLConfiguration(config Configuration) mysql.Configuration {
	return config.MySQL
}

func LoggingConfiguration(config Configuration) logging.Configuration {
	return config.Logging
}

func DiscordConfiguration(config Configuration) *discord.Configuration {
	return &config.Discord
}

func S3Configuration(config Configuration) *s3.Configuration {
	return &config.S3
}

func StorageConfiguration(config Configuration) *soundstore.Configuration {
	return &config.Storage
}

func KubernetesConfiguration(config Configuration) *kubernetes.Configuration {
	return &config.Kubernetes
}
func SaverConfiguration(config Configuration) *saver.Configuration {
	return &config.Saver
}

func HttpConfiguration(config Configuration) http.Configuration {
	return config.Http
}

func RedisConfiguration(config Configuration) redis.Configuration {
	return config.Redis
}
func PrometheusConfiguration(config Configuration) prometheus.Configuration {
        return config.Prometheus
}

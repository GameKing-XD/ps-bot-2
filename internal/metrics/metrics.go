package metrics

import (
	pprom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/fx"
)

type MetricsCollector struct {
	Prom          *prometheus.Prometheus
	soundCounters map[string]pprom.Counter
        SoundStore *soundstore.SoundStore
}

type NewMetricsCollectorParams struct {
        fx.In

	Prometheus *prometheus.Prometheus
        SoundStore *soundstore.SoundStore
}

func NewMetricsCollector(params NewMetricsCollectorParams) *MetricsCollector {
	return &MetricsCollector{
		Prom:          params.Prometheus,
		soundCounters: make(map[string]pprom.Counter),
                SoundStore: params.SoundStore,
	}
}

func (m *MetricsCollector) RegisterGuild(GuildID string) error {
        sounds := m.SoundStore.List(GuildID)
        for i := range sounds {
                m.RegisterPlaySound(GuildID, sounds[i])
        }

	return nil

}
func (m *MetricsCollector) RegisterPlaySound(GuildID, SoundName string) (pprom.Counter, error) {
	key := GuildID + "/" + SoundName
        var err error
	if _, ok := m.soundCounters[key]; !ok {

		counter := promauto.NewCounter(pprom.CounterOpts{
			ConstLabels: map[string]string{
				"guild_id": GuildID,
				"sound":    SoundName,
			},
			Name: "played_sound",
		})
                err = m.Prom.Register(counter)
                if err == nil {
		        m.soundCounters[key] = counter
                }
	}

	return m.soundCounters[key], err
}

func (m *MetricsCollector) PlaySound(GuildID, SoundName string) {
	key := GuildID + "/" + SoundName

        if _, ok := m.soundCounters[key]; ok {
                m.soundCounters[key].Inc()
        }
}

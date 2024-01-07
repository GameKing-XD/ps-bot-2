package web

import (
	echo "github.com/labstack/echo/v4"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Web struct {
	Repo  *repositories.GuildRepository
	Queue *queues.SoundsQueue
	Store *soundstore.SoundStore
	Log   *zap.Logger
}

type NewWebQueueParams struct {
	fx.In
	Repo  *repositories.GuildRepository
	Queue *queues.SoundsQueue
	Store *soundstore.SoundStore
	Log   *zap.Logger
        Prometheus *prometheus.Prometheus
}

func NewWeb(p NewWebQueueParams) (*Web, error) {

	return &Web{
		Repo:  p.Repo,
		Queue: p.Queue,
		Store: p.Store,
		Log:   p.Log.Named("web"),
	}, nil
}

func (w *Web) ApiGroup() string {
	return ""
}

func (w *Web) Version() string {
	return ""
}

type PostPlaySoundBody struct {
	Guild string `json:"guild" validate:"required"`
	Sound string `json:"sound" validate:"required"`
}

func (w *Web) Handler(e *echo.Group) {
	e.GET("", func(c echo.Context) error {
		return c.Redirect(302, "/app/index.html")
	})

	e.Static("app", "web/assets")

	e.GET("api/guilds", func(c echo.Context) error {
		m := w.Repo.GetGuilds()
		return c.JSON(200, m)
	})

	e.GET("api/sounds/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(422, map[string]string{"message": "path param id is required"})
		}
		s := w.Store.List(id)
		return c.JSON(200, s)
	})

	e.POST("api/play", func(c echo.Context) error {
		body := new(PostPlaySoundBody)

		if err := c.Bind(body); err != nil {
			return c.JSON(422, err)
		}

		err := w.Queue.Append(
			body.Guild,
			body.Sound,
		)

		if err != nil {
			return c.JSON(500, err)
		}

		return c.String(200, "")
	})
}

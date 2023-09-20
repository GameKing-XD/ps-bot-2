package web

import (

	"github.com/labstack/echo/v4"
	"github.com/tvanriel/ps-bot-2/internal/queues"
	"github.com/tvanriel/ps-bot-2/internal/repositories"
	"github.com/tvanriel/ps-bot-2/internal/soundstore"
)

type Web struct {
	repo  *repositories.GuildRepository
	store *soundstore.SoundStore
        queue *queues.SoundsQueue
}

func NewWeb(repo *repositories.GuildRepository, queue *queues.SoundsQueue, store *soundstore.SoundStore) (*Web, error) {

	return &Web{
		repo:  repo,
		queue: queue,
		store: store,
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
		m := w.repo.GetGuilds()
		return c.JSON(200, m)
	})

	e.GET("api/sounds/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(422, map[string]string{"message": "path param id is required"})
		}
		s := w.store.List(id)
		return c.JSON(200, s)
	})

	e.POST("api/play", func(c echo.Context) error {
		body := new(PostPlaySoundBody)

		if err := c.Bind(body); err != nil {
			return c.JSON(422, err)
		}

                err := w.queue.Append(body.Guild, body.Sound)

		if err != nil {
			return c.JSON(500, err)
		}

		return c.String(200, "")
	})
}

package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.48.0 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Warn("alias is empty")

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("alias is required"))

			return
		}

		err := deleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Error("url not found", slog.String("alias", alias))

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("url not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("deleted url with alias", slog.String("alias", alias))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Mager556/url-shortener/internal/lib/logger/sl"
	"github.com/Mager556/url-shortener/internal/lib/response"
	"github.com/Mager556/url-shortener/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Response struct {
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, s URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		url, err := s.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found")

			render.JSON(w, r, response.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to get URL", sl.Err(err), slog.String("alias", alias))

			render.JSON(w, r, response.Error("internal error"))
			return
		}
		log.Info("Successfull receipt of URL")
		log.Info("Redirecting to url...")

		http.Redirect(w, r, url, http.StatusFound)
	}
}

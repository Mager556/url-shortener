package delete

import (
	"log/slog"
	"net/http"

	"github.com/Mager556/url-shortener/internal/lib/logger/sl"
	resp "github.com/Mager556/url-shortener/internal/lib/response"
	"github.com/Mager556/url-shortener/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, s URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode json", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		log.Info("Request body decoded")

		if err := validator.New().Struct(req); err != nil {
			log.Error("Error while validating request")

			errs := err.(validator.ValidationErrors)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(errs))
			return
		}

		log.Info("Request body successfully validated")

		err = s.DeleteURL(req.Alias)
		if err != nil {
			log.Error("Failed to delete url", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			if err == storage.ErrURLNotFound {
				render.JSON(w, r, resp.Error(storage.ErrURLNotFound.Error()))
			} else {
				render.JSON(w, r, resp.Error("Failed to delete url"))
			}
			return
		}

		log.Info("Url successfully deleted")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp.OK())
	}
}

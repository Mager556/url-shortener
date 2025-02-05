package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Mager556/url-shortener/internal/lib/logger/sl"
	"github.com/Mager556/url-shortener/internal/lib/random"
	resp "github.com/Mager556/url-shortener/internal/lib/response"
	"github.com/Mager556/url-shortener/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias"`
}

type Reponse struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.1 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

const aliasLength = 6

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErrs := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.ValidationError(validateErrs))
			return
		}

		log.Info("request is validated")

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("failed to save url: url already exists", slog.String("url", req.URL))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save url", sl.Err(err), slog.String("url", req.URL))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("URL added", slog.Int64("id", id))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, Reponse{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}

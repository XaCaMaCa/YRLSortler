package save

import (
	"errors"
	"libmods/internal/lib/api/response"
	"libmods/internal/lib/logger/sl"
	"libmods/internal/lib/random"
	"libmods/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

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
			log.Error("Ошибка декода request body", sl.Err(err))
			render.JSON(w, r, response.Error("Ошибка декода request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr, ok := err.(validator.ValidationErrors)
			if !ok {
				log.Error("invalid request", sl.Err(err))
				render.JSON(w, r, response.Error("invalid request"))
				return
			}

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}

}
func responseError(w http.ResponseWriter, r *http.Request, err error) {
	render.JSON(w, r, response.Error(err.Error()))
}
func responseValidationError(w http.ResponseWriter, r *http.Request, errs validator.ValidationErrors) {
	render.JSON(w, r, response.ValidationError(errs))
}
func responseURLExists(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.Error("url already exists"))
}
func responseURLNotFound(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.Error("url not found"))
}
func responseURLAdded(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}
func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
	})
}

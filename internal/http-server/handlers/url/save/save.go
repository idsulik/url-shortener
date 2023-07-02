package save

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/idsulik/url-shortener/internal/lib/api/response"
	"github.com/idsulik/url-shortener/internal/logger"
	"net/http"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28 --name UrlSaver
type UrlSaver interface {
	SaveUrl(alias, url string) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.28 --name AliasGenerator
type AliasGenerator interface {
	NewAlias(size int) string
}

func New(log *logger.Logger, aliasGenerator AliasGenerator, urlSaver UrlSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req Request
		err := render.DecodeJSON(request.Body, &req)

		if err != nil {
			log.Error(err.Error())
			render.Status(request, http.StatusBadRequest)
			render.JSON(writer, request, resp.NewErrorResponse(err))
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error(err.Error())
			render.Status(request, http.StatusBadRequest)
			render.JSON(
				writer,
				request,
				resp.NewValidationErrorResponse(err.(validator.ValidationErrors)),
			)
			return
		}

		urlAlias := req.Alias
		if urlAlias == "" {
			urlAlias = aliasGenerator.NewAlias(6)
		}

		_, err = urlSaver.SaveUrl(urlAlias, req.Url)

		if err != nil {
			log.Error(err.Error())
			render.Status(request, http.StatusBadRequest)
			render.JSON(writer, request, resp.NewErrorResponse(err))
			return
		}

		render.Status(request, http.StatusCreated)
		render.JSON(writer, request, Response{Response: resp.NewOkResponse(), Alias: urlAlias})
	}
}

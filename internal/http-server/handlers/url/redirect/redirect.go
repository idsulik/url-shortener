package redirect

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "github.com/idsulik/url-shortener/internal/lib/api/response"
	"github.com/idsulik/url-shortener/internal/logger"
	"github.com/idsulik/url-shortener/internal/storage"
	"net/http"
)

type Request struct {
	Alias string `json:"url" validate:"required,url"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28 --name UrlGetter
type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *logger.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		alias := chi.URLParam(request, "alias")

		if alias == "" {
			log.Error("Alias is empty")
			render.Status(request, http.StatusNotFound)
			render.JSON(writer, request, resp.NewErrorResponse(storage.ErrUrlNotFound))
			return
		}

		url, err := urlGetter.GetUrl(alias)

		if err != nil {
			log.Error(err.Error())
			render.Status(request, http.StatusNotFound)
			render.JSON(writer, request, resp.NewErrorResponse(err))
			return
		}

		http.Redirect(writer, request, url, http.StatusFound)
	}
}

package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/idsulik/url-shortener/internal/http-server/handlers/url/save"
	"net/url"
	"testing"
)

const host = "localhost:8080"

func TestUrlShortener_HappyPath(t *testing.T) {
	u := url.URL{Scheme: "http", Host: host}

	e := httpexpect.Default(t, u.String())

	e.POST("/api/url/shorten").
		WithJSON(save.Request{
			Alias: gofakeit.Word(),
			Url:   gofakeit.URL(),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(201).
		JSON().
		Object().
		ContainsKey("alias")
}

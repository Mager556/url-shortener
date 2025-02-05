package url_shortener_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Mager556/url-shortener/internal/lib/api"
	"github.com/Mager556/url-shortener/internal/lib/random"
	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const host = "localhost:8080"

func Test_URLShortener_Save(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	for i := 0; i < 10; i++ {
		e.POST("/url").
			WithJSON(save.Request{
				URL:   gofakeit.URL(),
				Alias: random.NewRandomString(10),
			}).
			WithBasicAuth("myuser", "mypass").
			Expect().
			Status(200).
			JSON().Object().
			ContainsKey("alias")
		time.Sleep(time.Second / 8)
	}
}

func Test_URLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid url",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}
	for _, test := range testCases {
		u := url.URL{
			Scheme: "http",
			Host:   host,
		}
		e := httpexpect.Default(t, u.String())
		resp := e.POST("/url").
			WithJSON(save.Request{
				URL:   test.url,
				Alias: test.alias,
			}).
			WithBasicAuth("myuser", "mypass").
			Expect().
			JSON().Object()

		if test.error != "" {
			resp.NotContainsKey("alias")
			resp.Value("error").String().IsEqual(test.error)
			return
		}
		resp.Value("alias").String().IsEqual(test.alias)
		testRedirect(t, test.alias, test.url)
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToUrl, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, redirectedToUrl, urlToRedirect)
}

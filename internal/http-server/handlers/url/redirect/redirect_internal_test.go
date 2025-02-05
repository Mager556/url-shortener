package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/redirect"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"github.com/Mager556/url-shortener/internal/lib/api"
	"github.com/Mager556/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RedirectHandler(t *testing.T) {
	cases := []struct {
		name  string
		alias string
		url   string
	}{
		{
			name:  "Successfull",
			alias: "some_alias",
			url:   "httpls://google.com",
		},
	}

	for _, test := range cases {
		log := slogdiscard.NewDiscardLogger()
		urlGetter := mocks.NewURLGetter(t)

		if test.alias != "" {
			urlGetter.On("GetURL", test.alias).Return(test.url, nil).Once()
		}

		r := chi.NewRouter()
		r.Get("/{alias}", redirect.New(log, urlGetter))

		ts := httptest.NewServer(r)
		defer ts.Close()

		redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + test.alias)
		require.NoError(t, err)

		assert.Equal(t, test.url, redirectedToUrl)
	}
}

package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/save"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/save/mocks"
	"github.com/Mager556/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_SaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		url       string
		alias     string
		respError string
	}{
		{
			name:  "Success",
			url:   "https://www.google.com/",
			alias: "myalias",
		},

		{
			name: "Empty alias",
			url:  "https://www.google.com/",
		},

		{
			name:      "Empty url",
			alias:     "some_alias",
			respError: "field URL is required field",
		},
		{
			name:      "Invalid url",
			url:       "some invalid url",
			alias:     "alias",
			respError: "field URL is not a valid url",
		},
	}

	for _, test := range cases {
		urlSaver := mocks.NewURLSaver(t)

		if test.respError == "" {
			urlSaver.On("SaveURL", test.url, mock.AnythingOfType("string")).Return(int64(1), nil).Once()
		}

		discardLogger := slogdiscard.NewDiscardLogger()

		handler := save.New(discardLogger, urlSaver)

		input := fmt.Sprintf(`{
			"url": "%s",
			"alias": "%s"
		}`, test.url, test.alias)

		req := httptest.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		//assert.Equal(t, http.StatusOK, rec.Code)

		var resp save.Reponse

		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, test.respError, resp.Response.Error)
	}
}

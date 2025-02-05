package delete_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/delete"
	"github.com/Mager556/url-shortener/internal/http-server/handlers/url/delete/mocks"
	"github.com/Mager556/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/Mager556/url-shortener/internal/lib/response"
	"github.com/Mager556/url-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Sucsessfull",
			alias: "short",
		},
		{
			name:      "Empty alias",
			respError: "field Alias is required field",
		},
		{
			name:      "Not exist alias",
			alias:     "this_not_exists",
			respError: "url not found",
			mockError: storage.ErrURLNotFound,
		},
	}

	for _, test := range cases {
		URLDeleter := mocks.NewURLDeleter(t)
		log := slogdiscard.NewDiscardLogger()

		if test.respError == "" || test.mockError != nil {
			URLDeleter.On("DeleteURL", test.alias).Return(test.mockError).Once()
		}

		input := fmt.Sprintf(`{
			"alias": "%s"
		}`, test.alias)

		req := httptest.NewRequest(http.MethodDelete, "/", bytes.NewReader([]byte(input)))
		rec := httptest.NewRecorder()

		handler := delete.New(log, URLDeleter)

		handler.ServeHTTP(rec, req)

		var resp response.Response

		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, test.respError, resp.Error)
	}
}

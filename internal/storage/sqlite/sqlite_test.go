package sqlite_test

import (
	"testing"

	"github.com/Mager556/url-shortener/internal/storage"
	"github.com/Mager556/url-shortener/internal/storage/sqlite"
	"github.com/stretchr/testify/assert"
)

type TestParametr struct {
	UrlToSave     string
	Alias         string
	ExpectedError error
}

func Test_SaveAndDeleteUrl(t *testing.T) {
	tests := []TestParametr{

		{ // Default normal test
			UrlToSave:     "youtube.com",
			Alias:         "myshortref",
			ExpectedError: nil,
		},

		{ // Exists alias
			UrlToSave:     "mybadURL",
			Alias:         "myshortref",
			ExpectedError: storage.ErrURLExists,
		},
	}

	s, err := sqlite.New("./test_storages/test.db")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests { // deleting test alias if exist
		err := s.DeleteAlias(test.Alias)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, test := range tests {
		_, err = s.SaveURL(test.UrlToSave, test.Alias)
		assert.Equal(t, test.ExpectedError, err)
	}
}

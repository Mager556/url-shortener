package random_test

import (
	"testing"

	"github.com/Mager556/url-shortener/internal/lib/random"
	"github.com/stretchr/testify/assert"
)

func Test_NewRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "size = 1",
			length: 1,
		},

		{
			name:   "size = 3",
			length: 3,
		},

		{
			name:   "size = 5",
			length: 5,
		},

		{
			name:   "size = 10",
			length: 10,
		},

		{
			name:   "size = 20",
			length: 20,
		},
	}

	for _, value := range tests {
		s := random.NewRandomString(value.length)
		assert.Equal(t, value.length, len(s))
	}
}

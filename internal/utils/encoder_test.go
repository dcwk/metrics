package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		Data     string
		Expected string
	}{
		{
			Data:     "test",
			Expected: "test",
		},
	}

	for _, test := range tests {
		encodedData, err := Encrypt([]byte(test.Data), "./keys/public.pem")
		assert.NoError(t, err)

		decodedData, err := Decrypt(encodedData, "./keys/private.pem")
		assert.NoError(t, err)
		assert.Equal(t, test.Expected, string(decodedData))
	}
}

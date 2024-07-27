package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoder(t *testing.T) {
	tests := []struct {
		PublicKeyPath   string
		PrivateKeyPath  string
		Data            string
		ExpectedData    string
		PublicPemError  bool
		PrivatePemError bool
	}{
		{
			PublicKeyPath:   "./keys/public.pem",
			PrivateKeyPath:  "./keys/private.pem",
			Data:            "test",
			ExpectedData:    "test",
			PublicPemError:  false,
			PrivatePemError: false,
		},
		{
			PublicKeyPath:   "./keys/public.pem",
			PrivateKeyPath:  "./keys/private.pem",
			Data:            "",
			ExpectedData:    "",
			PublicPemError:  false,
			PrivatePemError: false,
		},
		{
			PublicKeyPath:   "./keys/public.pem",
			PrivateKeyPath:  "./keys/private.pem",
			Data:            "123123890###ololo",
			ExpectedData:    "123123890###ololo",
			PublicPemError:  false,
			PrivatePemError: false,
		},
		{
			PublicKeyPath:   "./keys/oloo.pem",
			PrivateKeyPath:  "./keys/private.pem",
			Data:            "123123890###ololo",
			ExpectedData:    "123123890###ololo",
			PublicPemError:  true,
			PrivatePemError: false,
		},
		{
			PublicKeyPath:   "./keys/public.pem",
			PrivateKeyPath:  "./keys/ololo.pem",
			Data:            "123123890###ololo",
			ExpectedData:    "123123890###ololo",
			PublicPemError:  false,
			PrivatePemError: true,
		},
	}

	for _, test := range tests {
		encodedData, err := Encrypt([]byte(test.Data), test.PublicKeyPath)
		if test.PublicPemError {
			assert.Error(t, err)
			assert.Nil(t, encodedData)
			continue
		} else {
			assert.NoError(t, err)
		}

		decodedData, err := Decrypt(encodedData, test.PrivateKeyPath)
		if test.PrivatePemError {
			assert.Error(t, err)
			assert.Nil(t, decodedData)
			continue
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedData, string(decodedData))
		}
	}
}

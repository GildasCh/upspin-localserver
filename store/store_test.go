package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"upspin.io/upspin"
)

func TestDial(t *testing.T) {
	store := Store{
		Debug: true,
	}

	actual, err := store.Dial(nil, upspin.Endpoint{})

	assert.NoError(t, err)
	assert.Equal(t, &store, actual)
}

func TestEndpoint(t *testing.T) {
	store := Store{
		Debug: true,
	}

	actual := store.Endpoint()

	assert.Equal(t, upspin.Endpoint{}, actual)
}

func TestClose(t *testing.T) {
	store := Store{
		Debug: true,
	}

	store.Close()
}

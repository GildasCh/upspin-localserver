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

func TestGetOK(t *testing.T) {
	store := Store{
		Root:  "../dir/test_data",
		Debug: true,
	}

	b, r, l, err := store.Get("abc-0")

	expected := []byte(`hello world!
`)

	assert.NoError(t, err)
	assert.Equal(t, expected, b)
	assert.Equal(t, &upspin.Refdata{Reference: "abc-0"}, r)
	assert.Equal(t, []upspin.Location(nil), l)
}

func TestGetHTTPBaseMetadataReturnsNotExist(t *testing.T) {
	_, _, _, err := (&Store{}).Get(upspin.HTTPBaseMetadata)

	assert.EqualError(t, err, "item does not exist")
}

func TestGetInvalidRefReturnsNotExist(t *testing.T) {
	_, _, _, err := (&Store{}).Get("something-notanumber")

	assert.EqualError(t, err, "item does not exist")
}

func TestGetErrorOpeningFileReturnsNotExist(t *testing.T) {
	_, _, _, err := (&Store{}).Get("missingfile-0")

	assert.EqualError(t, err, "item does not exist")
}

func TestGetUnreadableFileReturnsIOError(t *testing.T) {
	store := Store{
		Root:  "../dir/test_data",
		Debug: true,
	}

	_, _, _, err := store.Get(".-1048576")

	assert.EqualError(t, err, "I/O error")
}

func TestSplit(t *testing.T) {
	cases := map[string]struct {
		relativePath string
		offset       int64
	}{
		"filename.ext-0":                {"filename.ext", 0},
		"/in/subfolder/filename.ext-0":  {"/in/subfolder/filename.ext", 0},
		"/dir/filename-with-dash.ext-0": {"/dir/filename-with-dash.ext", 0},
	}

	store := Store{
		Root:  "../dir/test_data",
		Debug: true,
	}

	for in, expected := range cases {
		actualRelative, actualOffset, err := store.split(in)
		assert.NoError(t, err)
		assert.Equal(t, expected.relativePath, actualRelative)
		assert.Equal(t, expected.offset, actualOffset)
	}
}

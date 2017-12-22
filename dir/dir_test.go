package dir

import (
	"errors"
	"math/big"
	"testing"

	"github.com/gildasch/upspin-localserver/local"
	"github.com/gildasch/upspin-localserver/packing"
	"github.com/stretchr/testify/assert"
	_ "upspin.io/store/transports"
	"upspin.io/upspin"
)

type MockStorage struct {
	err error
}

func (ms *MockStorage) Stat(name string) (local.FileInfo, error) {
	if ms.err != nil {
		return local.FileInfo{}, ms.err
	}

	return local.FileInfo{
		Filename: "/test_data/abc",
		Dir:      "/test_data",
		IsDir:    false,
		Size:     20,
	}, nil
}

func (ms *MockStorage) List(pattern string) ([]local.FileInfo, error) {
	return nil, nil
}

type MockFactotum struct{}

func (mf *MockFactotum) FileSign(hash upspin.DEHash) (upspin.Signature, error) {
	R, S := big.NewInt(32948748), big.NewInt(982238482482302)
	return upspin.Signature{R, S}, nil
}

func (mf *MockFactotum) DirEntryHash(
	n, l upspin.PathName, a upspin.Attribute, p upspin.Packing,
	t upspin.Time, dkey, hash []byte) upspin.DEHash {
	return nil
}

type MockPacking struct{}

func (mp *MockPacking) DirEntry(username string, fi local.FileInfo, factotum packing.Factotum) *upspin.DirEntry {
	return &upspin.DirEntry{
		Sequence: 1234,
	}
}

func TestLookupOK(t *testing.T) {
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  &MockStorage{},
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},
	}

	entry, err := dir.Lookup("test.user@some-mail.com/test_data/abc")

	expected := &upspin.DirEntry{
		Sequence: 1234,
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, entry)
}

func TestLookupErrors(t *testing.T) {
	storage := &MockStorage{}
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  storage,
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},
	}

	_, err := dir.Lookup("test.usersome-mail.com/test_data/abc")
	assert.EqualError(t, err, "error parsing path: user.Parse: user test.usersome-mail.com: invalid operation: user name must contain one @ symbol")
	_, err = dir.Lookup("user.test@some-mail.com/test_data/abc")
	assert.EqualError(t, err, "user \"user.test@some-mail.com\" is not known on this server")

	storage.err = errors.New("dummy error")
	_, err = dir.Lookup("test.user@some-mail.com/test_data/abc")
	assert.EqualError(t, err, "could not stat file \"test_data/abc\"")
}

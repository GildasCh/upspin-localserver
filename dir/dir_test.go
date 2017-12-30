package dir

import (
	"errors"
	"math/big"
	"testing"

	"github.com/gildasch/upspin-localserver/local"
	"github.com/gildasch/upspin-localserver/packing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"upspin.io/access"
	"upspin.io/config"
	_ "upspin.io/store/transports"
	"upspin.io/upspin"
)

type MockStorage struct {
	err       error
	statError error
	listError error
}

func (ms *MockStorage) Stat(name string) (local.FileInfo, error) {
	if ms.err != nil {
		return local.FileInfo{}, ms.err
	}

	if ms.statError != nil {
		return local.FileInfo{}, ms.statError
	}

	if name == "/test_data" {
		return local.FileInfo{
			Filename: "/test_data",
			Dir:      "/",
			IsDir:    true,
		}, nil
	}

	return local.FileInfo{
		Filename: "/test_data/abc",
		Dir:      "/test_data",
		IsDir:    false,
		Size:     20,
	}, nil
}

func (ms *MockStorage) List(pattern string) ([]local.FileInfo, error) {
	if ms.err != nil {
		return nil, ms.err
	}

	if ms.listError != nil {
		return nil, ms.listError
	}

	return []local.FileInfo{
		local.FileInfo{
			Filename: "/test_data/abc",
			Dir:      "/test_data",
			IsDir:    false,
			Size:     20,
		},
		local.FileInfo{
			Filename: "/test_data/cba",
			Dir:      "/test_data",
			IsDir:    false,
			Size:     40,
		}}, nil
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
	if fi.Filename == "/test_data/cba" {
		return &upspin.DirEntry{
			Name:     upspin.PathName(username + fi.Filename),
			Sequence: 4321,
		}
	}
	return &upspin.DirEntry{
		Name:     upspin.PathName(username + fi.Filename),
		Sequence: 1234,
	}
}

func TestDial(t *testing.T) {
	userName := "test.user@some-mail.com"
	defaultAccess, err := access.New(upspin.PathName(userName + "/"))
	require.NoError(t, err)
	dir := Dir{
		Username:      userName,
		Root:          ".",
		Storage:       &MockStorage{},
		Debug:         false,
		Factotum:      &MockFactotum{},
		Packing:       &MockPacking{},
		defaultAccess: defaultAccess,
	}

	endpoint := upspin.Endpoint{
		Transport: upspin.InProcess,
		NetAddr:   "", // ignored
	}
	cfg := config.New()
	cfg = config.SetUserName(cfg, upspin.UserName(userName))
	cfg = config.SetPacking(cfg, upspin.EEPack)
	cfg = config.SetKeyEndpoint(cfg, endpoint)
	cfg = config.SetStoreEndpoint(cfg, endpoint)
	cfg = config.SetDirEndpoint(cfg, endpoint)

	actualService, err := dir.Dial(cfg, upspin.Endpoint{})
	require.NoError(t, err)

	actual, ok := actualService.(*Dir)
	require.True(t, ok)
	assert.Equal(t, cfg.UserName(), actual.userName)
	assert.True(t, actual.dialed)
}

func TestEndpoint(t *testing.T) {
	dir := Dir{
		Debug: false,
	}

	actual := dir.Endpoint()

	assert.Equal(t, upspin.Endpoint{}, actual)
}

func TestClose(t *testing.T) {
	dir := Dir{
		Debug: false,
	}

	dir.Close()
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
		Name:     "test.user@some-mail.com/test_data/abc",
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
	assert.EqualError(t, err, "could not stat file \"test_data/abc\": dummy error")
	storage.err = nil
}

func TestGlobOK(t *testing.T) {
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  &MockStorage{},
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},
	}

	entries, err := dir.Glob("test.user@some-mail.com/test_data/*")

	expected := []*upspin.DirEntry{
		&upspin.DirEntry{
			Name:     "test.user@some-mail.com/test_data/abc",
			Sequence: 1234},
		&upspin.DirEntry{
			Name:     "test.user@some-mail.com/test_data/cba",
			Sequence: 4321}}

	assert.NoError(t, err)
	assert.Equal(t, expected, entries)
}

func TestGlobErrors(t *testing.T) {
	storage := &MockStorage{}
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  storage,
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},
	}

	_, err := dir.Glob("user.test@some-mail.com/test_data/*")
	assert.EqualError(t, err, "path unknown")

	storage.listError = errors.New("dummy error")
	_, err = dir.Glob("test.user@some-mail.com/test_data/*")
	assert.EqualError(t, err, "error during glob: test.user@some-mail.com/test_data: error reading dir: dummy error")
}

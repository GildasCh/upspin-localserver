package dir

import (
	"math/big"
	"testing"

	"github.com/gildasch/upspin-localserver/local"
	"github.com/stretchr/testify/assert"
	_ "upspin.io/store/transports"
	"upspin.io/upspin"
)

type MockStorage struct{}

func (ms *MockStorage) Stat(name string) (local.FileInfo, error) {
	return local.FileInfo{
		Filename: "/test_data/abc",
		Dir:      ".",
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

func TestLookup(t *testing.T) {
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  &MockStorage{},
		Debug:    false,
		Factotum: &MockFactotum{},
	}

	entry, err := dir.Lookup("test.user@some-mail.com/test_data/abc")

	expected := &upspin.DirEntry{
		Packing: upspin.PlainPack,
		Blocks: []upspin.DirBlock{
			upspin.DirBlock{
				Location: upspin.Location{
					Endpoint: upspin.Endpoint{
						Transport: upspin.Remote,
						NetAddr:   "usl.gildas.ch"},
					Reference: "/test_data/abc-0"},
				Offset:   0,
				Size:     20,
				Packdata: []uint8(nil)}},
		Packdata: []uint8{0x8, 0x1, 0xf6, 0xc2, 0xc, 0xe, 0x3, 0x7d, 0x57, 0x37, 0xd5, 0x44, 0x7e, 0x0, 0x0},
		Writer:   "test.user@some-mail.com",
		Name:     "test.user@some-mail.com/test_data/abc",
		Sequence: 0,
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, entry)
}

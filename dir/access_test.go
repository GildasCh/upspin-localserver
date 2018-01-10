package dir

import (
	"testing"

	"github.com/gildasch/upspin-localserver/local"
	"github.com/stretchr/testify/assert"
)

type MockAccessStorage struct {
	accessed map[string]string
	called   []string
}

func (ms *MockAccessStorage) Stat(name string) (local.FileInfo, error) {
	return local.FileInfo{}, nil
}

func (ms *MockAccessStorage) List(pattern string) ([]local.FileInfo, error) {
	return nil, nil
}

func (ms *MockAccessStorage) Access(name string) ([]byte, bool) {
	ms.called = append(ms.called, name)
	found, ok := ms.accessed[name]
	if !ok {
		return nil, false
	}
	return []byte(found), ok
}

func TestCanList(t *testing.T) {
	storage := &MockAccessStorage{}
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  storage,
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},

		userName: "test.user@some-mail.com",
	}

	storage.accessed = map[string]string{
		"target.user@mail.com/a/Access": "*:test.user@some-mail.com",
	}
	actual := dir.canList("target.user@mail.com/a/dir/somewhere/about")

	assert.Equal(t, true, actual)

	expectedCalled := []string{
		"target.user@mail.com/a/dir/somewhere/about/Access",
		"target.user@mail.com/a/dir/somewhere/Access",
		"target.user@mail.com/a/dir/Access",
		"target.user@mail.com/a/Access"}

	assert.Equal(t, expectedCalled, storage.called)
}

func TestCanListNoAccessFile(t *testing.T) {
	storage := &MockAccessStorage{}
	dir := Dir{
		Username: "test.user@some-mail.com",
		Root:     ".",
		Storage:  storage,
		Debug:    false,
		Factotum: &MockFactotum{},
		Packing:  &MockPacking{},

		userName: "test.user@some-mail.com",
	}

	actual := dir.canList("target.user@mail.com/a/dir/somewhere/about")

	assert.Equal(t, false, actual)

	expectedCalled := []string{
		"target.user@mail.com/a/dir/somewhere/about/Access",
		"target.user@mail.com/a/dir/somewhere/Access",
		"target.user@mail.com/a/dir/Access",
		"target.user@mail.com/a/Access",
		"target.user@mail.com/Access"}

	assert.Equal(t, expectedCalled, storage.called)
}

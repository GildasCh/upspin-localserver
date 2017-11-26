package dir

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"upspin.io/upspin"
)

//go:generate mkdir test_data/empty

func Test_Dial(t *testing.T) {
	d := &Dir{}

	actual, err := d.Dial(nil, upspin.Endpoint{})

	assert.NoError(t, err)
	assert.Equal(t, d, actual)
}

func Test_Endpoint(t *testing.T) {
	d := &Dir{}

	expected := upspin.Endpoint{}

	assert.Equal(t, expected, d.Endpoint())
}

func Test_Lookup_OK(t *testing.T) {
	d := &Dir{Username: "test@gmail.com"}

	de, err := d.Lookup("test@gmail.com/toto")

	expected := &upspin.DirEntry{}

	assert.NoError(t, err)
	assert.Equal(t, expected, de)
}

func Test_Lookup_Invalid_Path(t *testing.T) {
	d := &Dir{Username: "test@gmail.com"}

	de, err := d.Lookup("invalid-path")

	assert.Error(t, err)
	assert.Nil(t, de)
}

func Test_Lookup_Different_Username(t *testing.T) {
	d := &Dir{Username: "test@gmail.com"}

	de, err := d.Lookup("test2@gmail.com/toto")

	assert.Error(t, err)
	assert.Nil(t, de)
}

func Test_Glob_OK(t *testing.T) {
	d := &Dir{
		Username: "test@gmail.com",
		Root:     "test_data"}

	des, err := d.Glob("test@gmail.com")

	expected := []*upspin.DirEntry{
		&upspin.DirEntry{
			Blocks: []upspin.DirBlock{upspin.DirBlock{
				Location: upspin.Location{
					Reference: "caca"},
				Size: 13}},
			Name: "abc"},
		&upspin.DirEntry{
			Blocks: []upspin.DirBlock{upspin.DirBlock{
				Location: upspin.Location{
					Reference: "caca"}}},
			Name: "cde"},
		&upspin.DirEntry{
			Attr: upspin.AttrDirectory,
			Name: "empty"},
		&upspin.DirEntry{
			Attr: upspin.AttrDirectory,
			Name: "subdir"},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, des)
}

func Test_Glob_OK_Subpath(t *testing.T) {
	d := &Dir{
		Username: "test@gmail.com",
		Root:     "test_data"}

	des, err := d.Glob("test@gmail.com/subdir")

	expected := []*upspin.DirEntry{
		&upspin.DirEntry{
			Blocks: []upspin.DirBlock{upspin.DirBlock{
				Location: upspin.Location{
					Reference: "caca"},
				Size: 16}},
			Name: "fgh"},
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, des)
}

func Test_Glob_OK_Empty_Subpath(t *testing.T) {
	d := &Dir{
		Username: "test@gmail.com",
		Root:     "test_data"}

	des, err := d.Glob("test@gmail.com/empty")

	expected := []*upspin.DirEntry{}

	assert.NoError(t, err)
	assert.Equal(t, expected, des)
}

func Test_Glob_Wrong_Username(t *testing.T) {
	d := &Dir{
		Username: "test@gmail.com",
		Root:     "test_data"}

	des, err := d.Glob("testg@mail.com")

	assert.Error(t, err)
	assert.Nil(t, des)
}

func Test_Glob_Wrong_Dir_Root(t *testing.T) {
	d := &Dir{
		Username: "test@gmail.com",
		Root:     "wrong_root"}

	des, err := d.Glob("test@gmail.com")

	assert.Error(t, err)
	assert.Nil(t, des)
}

package local

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenOK(t *testing.T) {
	cases := map[string][]byte{
		"test_1.txt": []byte(`some text...
...
`),
		"subdir/../../../test_1.txt": []byte(`some text...
...
`),
		"unknown_dir/../../../test_1.txt": []byte(`some text...
...
`),
		"subdir/toto": []byte(`
`),
	}

	s := &Storage{"test_data"}

	for in, expected := range cases {
		f, err := s.Open(in)
		assert.NoError(t, err)
		bytes, err := ioutil.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, expected, bytes)
	}
}

func TestOpenNotFound(t *testing.T) {
	cases := []string{
		"test_2.txt",
		"unknown_dir/test_1.txt",
	}

	s := &Storage{"test_data"}

	for _, in := range cases {
		f, err := s.Open(in)
		assert.Error(t, err)
		assert.Nil(t, f)
	}
}

func TestStatOK(t *testing.T) {
	fi1 := &FileInfo{
		Filename: "test_data/test_1.txt",
		Dir:      "test_data",
		IsDir:    false,
		Size:     17}
	fi2 := &FileInfo{
		Filename: "test_data/subdir/toto",
		Dir:      "test_data/subdir",
		IsDir:    false,
		Size:     1}

	cases := map[string]*FileInfo{
		"test_1.txt":                      fi1,
		"subdir/../../../test_1.txt":      fi1,
		"unknown_dir/../../../test_1.txt": fi1,
		"subdir/toto":                     fi2,
	}

	s := &Storage{"test_data"}

	for in, expected := range cases {
		fi, err := s.Stat(in)
		assert.NoError(t, err)
		assert.Equal(t, expected, fi)
	}
}

func TestStatNotFound(t *testing.T) {
	cases := []string{
		"test_2.txt",
		"unknown_dir/test_1.txt",
	}

	s := &Storage{"test_data"}

	for _, in := range cases {
		f, err := s.Stat(in)
		assert.Error(t, err)
		assert.Nil(t, f)
	}
}

func TestListOK(t *testing.T) {
	expected1 := []FileInfo{
		FileInfo{
			Filename: "subdir",
			Dir:      ".",
			IsDir:    true,
			Size:     4096},
		FileInfo{
			Filename: "test_1.txt",
			Dir:      ".",
			IsDir:    false,
			Size:     17}}
	expected2 := []FileInfo{
		FileInfo{
			Filename: "toto",
			Dir:      "test_data",
			IsDir:    false,
			Size:     1}}

	cases := map[string][]FileInfo{
		"/":      expected1,
		"subdir": expected2,
	}

	s := &Storage{"test_data"}

	for in, expected := range cases {
		fis, err := s.List(in)
		assert.NoError(t, err)
		assert.Equal(t, expected, fis)
	}
}

func TestListNotFound(t *testing.T) {
	cases := []string{
		"subdir/subdir",
		"unknown_dir",
		"subdir/toto",
	}

	s := &Storage{"test_data"}

	for _, in := range cases {
		fis, err := s.List(in)
		assert.Error(t, err)
		assert.Equal(t, []FileInfo(nil), fis)
	}
}

func TestFilename(t *testing.T) {
	cases := map[string]string{
		"toto":                  "a/toto",
		"a/toto":                "a/a/toto",
		"../b/toto":             "a/b/toto",
		"../b/../../../../toto": "a/toto",
	}

	s := &Storage{"a"}

	for in, expected := range cases {
		assert.Equal(t, expected, s.filename(in))
	}
}

func TestDir(t *testing.T) {
	cases := map[string]string{
		"toto":                  "a",
		"a/toto":                "a/a",
		"../b/toto":             "a/b",
		"../b/../../../../toto": "a",
	}

	s := &Storage{"a"}

	for in, expected := range cases {
		assert.Equal(t, expected, s.dir(in))
	}
}

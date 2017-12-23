package local

import (
	"io/ioutil"
	"testing"
	"time"

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

func TestDeeperRoot(t *testing.T) {
	expected := []byte(`
`)

	s := &Storage{"test_data/subdir"}

	f, err := s.Open("toto")
	assert.NoError(t, err)
	bytes, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, expected, bytes)

	assert.Equal(t, "test_data/subdir/toto", s.filename("toto"))
	assert.Equal(t, "test_data/subdir", s.dir("toto"))
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
	t1, _ := time.Parse(
		"2006-01-02 15:04:05.000000000 -0700 MST",
		"2017-12-21 00:30:01.020556451 +0100 CET")
	fi1 := FileInfo{
		Filename: "/test_1.txt",
		Dir:      "",
		IsDir:    false,
		Size:     17,
		Time:     t1}
	t2, _ := time.Parse(
		"2006-01-02 15:04:05.000000000 -0700 MST",
		"2017-12-21 00:30:17.880470714 +0100 CET")
	fi2 := FileInfo{
		Filename: "/subdir/toto",
		Dir:      "/subdir",
		IsDir:    false,
		Size:     1,
		Time:     t2}

	cases := map[string]FileInfo{
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
		assert.Zero(t, f)
	}
}

func TestListOK(t *testing.T) {
	t1, _ := time.Parse(
		"2006-01-02 15:04:05.000000000 -0700 MST",
		"2017-12-21 00:30:17.887470679 +0100 CET")
	t2, _ := time.Parse(
		"2006-01-02 15:04:05.000000000 -0700 MST",
		"2017-12-21 00:30:01.020556451 +0100 CET")
	expected1 := []FileInfo{
		FileInfo{
			Filename: "/subdir",
			Dir:      ".",
			IsDir:    true,
			Size:     4096,
			Time:     t1},
		FileInfo{
			Filename: "/test_1.txt",
			Dir:      ".",
			IsDir:    false,
			Size:     17,
			Time:     t2}}
	t3, _ := time.Parse(
		"2006-01-02 15:04:05.000000000 -0700 MST",
		"2017-12-21 00:30:17.880470714 +0100 CET")
	expected2 := []FileInfo{
		FileInfo{
			Filename: "subdir/toto",
			Dir:      "test_data",
			IsDir:    false,
			Size:     1,
			Time:     t3}}

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

package local

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Storage struct {
	Root string
}

type FileInfo struct {
	Filename string
	Dir      string
	IsDir    bool
	Size     int64
	Time     time.Time
}

func (fi FileInfo) Path() string {
	return filepath.Join(fi.Dir, fi.Filename)
}

func (s *Storage) Open(name string) (*os.File, error) {
	f, err := os.Open(s.filename(name))
	if err != nil {
		return nil, errors.Wrapf(err, "could not open file %q", name)
	}

	return f, nil
}

func (s *Storage) Stat(name string) (FileInfo, error) {
	f, err := os.Open(s.filename(name))
	if err != nil {
		return FileInfo{}, errors.Wrapf(err, "could not open file %q", name)
	}
	fi, err := f.Stat()
	if err != nil {
		return FileInfo{}, errors.Wrapf(err, "could not stat file %q", name)
	}

	return FileInfo{
		Filename: strings.TrimPrefix(s.filename(name), s.Root),
		Dir:      strings.TrimPrefix(s.dir(name), s.Root),
		IsDir:    fi.IsDir(),
		Size:     fi.Size(),
		Time:     fi.ModTime(),
	}, nil
}

func (s *Storage) Access(name string) ([]byte, bool) {
	return []byte(`read: gildaschbt@gmail.com`), true
}

func (s *Storage) List(pattern string) ([]FileInfo, error) {
	fis, err := ioutil.ReadDir(s.filename(pattern))
	if err != nil {
		return nil, errors.Wrap(err, "error reading dir")
	}

	infos := []FileInfo{}

	for _, fi := range fis {
		infos = append(infos, FileInfo{
			Filename: filepath.Join(pattern, fi.Name()),
			Dir:      s.dir(pattern),
			IsDir:    fi.IsDir(),
			Size:     fi.Size(),
			Time:     fi.ModTime(),
		})
	}

	return infos, nil
}

func (s *Storage) filename(name string) string {
	return filepath.Join(s.Root, filepath.FromSlash(filepath.Clean("/"+name)))
}

func (s *Storage) dir(name string) string {
	return filepath.Dir(s.filename(name))
}

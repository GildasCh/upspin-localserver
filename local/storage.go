package local

import (
	"io/ioutil"
	"os"
	"path/filepath"

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
		Filename: s.filename(name),
		Dir:      s.dir(name),
		IsDir:    fi.IsDir(),
		Size:     fi.Size(),
	}, nil
}

func (s *Storage) List(pattern string) ([]FileInfo, error) {
	fis, err := ioutil.ReadDir(s.filename(pattern))
	if err != nil {
		return nil, errors.Wrap(err, "error reading dir")
	}

	infos := []FileInfo{}

	for _, fi := range fis {
		infos = append(infos, FileInfo{
			Filename: fi.Name(),
			Dir:      s.dir(pattern),
			IsDir:    fi.IsDir(),
			Size:     fi.Size(),
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

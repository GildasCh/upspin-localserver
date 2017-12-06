package dir

import (
	"fmt"
	"io/ioutil"
	"os"
	gopath "path"
	"strings"

	"github.com/pkg/errors"
	"upspin.io/path"
	"upspin.io/upspin"
)

type Dir struct {
	upspin.DirServer

	Username string

	Root string
}

func (d *Dir) Dial(upspin.Config, upspin.Endpoint) (upspin.Service, error) {
	return d, nil
}

func (d *Dir) Endpoint() upspin.Endpoint {
	return upspin.Endpoint{}
}

func (d *Dir) Close() {}

func (d *Dir) Lookup(name upspin.PathName) (*upspin.DirEntry, error) {
	p, err := path.Parse(name)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing path")
	}
	if string(p.User()) != d.Username {
		return nil,
			fmt.Errorf("user %q is not known on this server", p.User())
	}

	return &upspin.DirEntry{}, nil
}

func (d *Dir) Glob(pattern string) ([]*upspin.DirEntry, error) {
	if !strings.HasPrefix(pattern, d.Username) {
		return nil, errors.New("path unknown")
	}
	pattern = strings.TrimPrefix(pattern, d.Username)
	pattern = strings.TrimSuffix(pattern, "*")
	localPath := gopath.Join(d.Root, pattern)
	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return nil, errors.Wrap(err, "error reading dir")
	}

	ret := []*upspin.DirEntry{}

	for _, f := range files {
		de := &upspin.DirEntry{
			Name: upspin.PathName(f.Name()),
		}
		if f.IsDir() {
			de.Attr = upspin.AttrDirectory
		} else {
			de.Blocks = blocksFromFile(f)
		}

		ret = append(ret, de)
	}

	return ret, nil
}

func blocksFromFile(f os.FileInfo) (dbs []upspin.DirBlock) {
	size := f.Size()
	offset := int64(0)
	for size > 0 {
		s := int64(upspin.BlockSize)
		if s > size {
			s = size
		}
		size -= s
		ref := fmt.Sprintf("%s-%d", f.Name(), offset)
		dbs = append(dbs, upspin.DirBlock{
			Location: upspin.Location{
				Endpoint: upspin.Endpoint{
					Transport: upspin.Remote,
					NetAddr:   "usl.gildas.ch",
				},
				Reference: upspin.Reference(ref)},
			Offset: offset,
			Size:   s,
		})
		offset += s
	}

	return
}

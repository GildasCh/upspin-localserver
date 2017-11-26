package dir

import (
	"fmt"

	"github.com/pkg/errors"
	"upspin.io/path"
	"upspin.io/upspin"
)

type Dir struct {
	upspin.DirServer

	username string
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
	if string(p.User()) != d.username {
		return nil,
			fmt.Errorf("user %q is not known on this server", p.User())
	}

	return &upspin.DirEntry{}, nil
}

func (d *Dir) Glob(pattern string) ([]*upspin.DirEntry, error) {
	return []*upspin.DirEntry{
		&upspin.DirEntry{Name: "abc"},
		&upspin.DirEntry{Name: "def"},
	}, nil
}

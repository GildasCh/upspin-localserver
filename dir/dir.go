package dir

import (
	"fmt"
	"io/ioutil"
	"os"
	gopath "path"
	"strings"

	"github.com/gildasch/upspin-localserver/packing"
	"github.com/pkg/errors"
	"upspin.io/path"
	"upspin.io/upspin"
)

type Dir struct {
	upspin.DirServer

	Username string
	Root     string
	Debug    bool
	Config   upspin.Config
}

func (d *Dir) Dial(config upspin.Config, endpoint upspin.Endpoint) (upspin.Service, error) {
	if d.Debug {
		fmt.Printf("dir.Dial called with config=%#v, endpoint=%#v\n", config, endpoint)
	}

	return d, nil
}

func (d *Dir) Endpoint() upspin.Endpoint {
	if d.Debug {
		fmt.Printf("dir.Endpoint called\n")
	}

	return upspin.Endpoint{}
}

func (d *Dir) Close() {
	if d.Debug {
		fmt.Printf("dir.Close called\n")
	}
}

func (d *Dir) Lookup(name upspin.PathName) (*upspin.DirEntry, error) {
	if d.Debug {
		fmt.Printf("dir.Lookup called with name=%#v\n", name)
	}

	p, err := path.Parse(name)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing path")
	}
	if string(p.User()) != d.Username {
		return nil,
			fmt.Errorf("user %q is not known on this server", p.User())
	}

	f, err := os.Open(gopath.Join(d.Root, p.FilePath()))
	if err != nil {
		return nil,
			fmt.Errorf("could not open file %q", p.FilePath())
	}
	fi, err := f.Stat()
	if err != nil {
		return nil,
			fmt.Errorf("could not stat file %q", p.FilePath())
	}

	de := packing.PlainDirEntry(fi, d.Config)

	if d.Debug {
		fmt.Printf("dir.Lookup returning %#v\n", de)
	}

	return de, nil
}

func (d *Dir) Glob(pattern string) ([]*upspin.DirEntry, error) {
	if d.Debug {
		fmt.Printf("dir.Glob called with pattern=%#v\n", pattern)
	}

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

	for _, fi := range files {
		de := packing.PlainDirEntry(fi, d.Config)
		ret = append(ret, de)
	}

	if d.Debug {
		fmt.Printf("dir.Glob returning %#v\n", ret)
	}

	return ret, nil
}

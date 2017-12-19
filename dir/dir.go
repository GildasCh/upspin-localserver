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

const packdata = "nothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothingnothing"

type Dir struct {
	upspin.DirServer

	Username string
	Root     string
	Debug    bool
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

	de := dirEntryFromFileInfo(fi)

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

	for _, f := range files {
		de := dirEntryFromFileInfo(f)
		ret = append(ret, de)
	}

	if d.Debug {
		fmt.Printf("dir.Glob returning %#v\n", ret)
	}

	return ret, nil
}

func dirEntryFromFileInfo(f os.FileInfo) *upspin.DirEntry {
	de := &upspin.DirEntry{
		Name:     upspin.PathName("gildaschbt+local@gmail.com/" + f.Name()),
		Packing:  upspin.PlainPack,
		Packdata: []byte(packdata),
	}
	if f.IsDir() {
		de.Attr = upspin.AttrDirectory
	} else {
		de.Blocks = blocksFromFileInfo(f)
	}
	return de
}

func blocksFromFileInfo(f os.FileInfo) (dbs []upspin.DirBlock) {
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
			Offset:   offset,
			Size:     s,
			Packdata: []byte(packdata),
		})
		offset += s
	}

	return
}

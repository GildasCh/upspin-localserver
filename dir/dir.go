package dir

import (
	"fmt"
	"strings"

	"github.com/gildasch/upspin-localserver/local"
	"github.com/gildasch/upspin-localserver/packing"
	"github.com/pkg/errors"
	"upspin.io/access"
	uperrors "upspin.io/errors"
	"upspin.io/path"
	"upspin.io/serverutil"
	"upspin.io/upspin"
	"upspin.io/user"
	"upspin.io/valid"
)

type Storage interface {
	Stat(name string) (local.FileInfo, error)
	List(pattern string) ([]local.FileInfo, error)
	Access(name string) []byte
}

type Dir struct {
	upspin.DirServer

	Username string
	Root     string
	Storage  Storage
	Debug    bool
	Factotum packing.Factotum
	Packing  packing.Simulator

	// userName is the name of the user on behalf of whom this
	// server is serving.
	userName upspin.UserName

	// baseUser, suffix and domain are the components of userName as parsed
	// by user.Parse.
	userBase, userSuffix, userDomain string

	// dialed reports whether the instance was created using Dial, not New.
	dialed bool

	// defaultAccess is the parsed empty Access files that implicitly exists
	// at the root of every user's tree, if an explicit one is not found.
	defaultAccess *access.Access
}

func (d *Dir) Dial(ctx upspin.Config, endpoint upspin.Endpoint) (upspin.Service, error) {
	if d.Debug {
		// fmt.Printf("dir.Dial called with ctx=%#v, endpoint=%#v, ctx.UserName()=%#v, ctx.Factotum()=%#v, ctx.Packing()=%#v, ctx.KeyEndpoint()=%#v, ctx.DirEndpoint()=%#v, ctx.StoreEndpoint()=%#v\n", ctx, endpoint, ctx.UserName(), ctx.Factotum(), ctx.Packing(), ctx.KeyEndpoint(), ctx.DirEndpoint(), ctx.StoreEndpoint())
		fmt.Printf("dir.Dial called with ctx=%#v, endpoint=%#v\n", ctx, endpoint)
	}

	if err := valid.UserName(ctx.UserName()); err != nil {
		return nil, errors.Wrapf(err, "invalid username")
	}

	cp := *d // copy of the generator instance.
	// Overwrite the userName and its sub-components (base, suffix, domain).
	cp.userName = ctx.UserName()
	cp.dialed = true
	var err error
	cp.userBase, cp.userSuffix, cp.userDomain, err = user.Parse(cp.userName)
	if err != nil {
		return nil, err
	}

	// create a default Access file for this user.
	cp.defaultAccess, err = access.New(upspin.PathName(cp.userName + "/"))
	if err != nil {
		return nil, err
	}
	return &cp, nil
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

	fi, err := d.Storage.Stat(p.FilePath())
	if err != nil {
		return nil,
			fmt.Errorf("could not stat file %q: %v", p.FilePath(), err)
	}

	de := d.Packing.DirEntry(d.Username, fi, d.Factotum)

	if d.Debug {
		fmt.Printf("dir.Lookup returning %#v\n", de)
	}

	return de, nil
}

func (d *Dir) Glob(pattern string) ([]*upspin.DirEntry, error) {
	if d.Debug {
		fmt.Printf("dir.Glob called with pattern=%#v, d=%#v\n", pattern, d)
	}

	if !strings.HasPrefix(pattern, d.Username) {
		return nil, errors.New("path unknown")
	}

	entries, err := serverutil.Glob(pattern, d.Lookup, d.listDir)
	if err != nil && err != upspin.ErrFollowLink {
		return nil, errors.Wrap(err, "error during glob")
	}

	if d.Debug {
		fmt.Printf("dir.Glob returning %#v\n", entries)
	}

	return entries, err
}

func (d *Dir) listDir(name upspin.PathName) ([]*upspin.DirEntry, error) {
	if d.Debug {
		fmt.Printf("dir.listDir called with name=%#v, d=%#v\n", name, d)
	}

	pattern := strings.TrimPrefix(string(name), d.Username)

	if d.userName == "" {
		return nil, uperrors.E(uperrors.Private)
	}

	acc, err := access.Parse(name+"/Access", d.Storage.Access(pattern+"Access"))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing access file")
	}

	if ok, _ := acc.Can(d.userName, access.List, upspin.PathName(pattern), nil); !ok {
		return nil, uperrors.E(uperrors.Private)
	}

	fis, err := d.Storage.List(pattern)
	if err != nil {
		return nil, errors.Wrap(err, "error reading dir")
	}

	ret := []*upspin.DirEntry{}

	for _, fi := range fis {
		de := d.Packing.DirEntry(d.Username, fi, d.Factotum)
		ret = append(ret, de)
	}

	if d.Debug {
		fmt.Printf("dir.listDir returning %#v\n", ret)
	}

	return ret, nil

}

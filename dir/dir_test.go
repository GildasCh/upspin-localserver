package dir

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"upspin.io/bind"
	"upspin.io/config"
	"upspin.io/dir/inprocess"
	"upspin.io/factotum"
	"upspin.io/path"
	_ "upspin.io/store/transports"
	"upspin.io/test/testutil"
	"upspin.io/upspin"
)

func newConfigAndServices(name upspin.UserName) (cfg upspin.Config, key upspin.KeyServer, dir upspin.DirServer, store upspin.StoreServer) {
	endpoint := upspin.Endpoint{
		Transport: upspin.InProcess,
		NetAddr:   "", // ignored
	}
	cfg = config.New()
	cfg = config.SetUserName(cfg, name)
	cfg = config.SetPacking(cfg, upspin.EEPack)
	cfg = config.SetKeyEndpoint(cfg, endpoint)
	cfg = config.SetStoreEndpoint(cfg, endpoint)
	cfg = config.SetDirEndpoint(cfg, endpoint)
	f, err := factotum.NewFromDir(testutil.Repo("key", "testdata", "user1")) // Always use user1's keys.
	if err != nil {
		panic(err)
	}
	cfg = config.SetFactotum(cfg, f)

	key, _ = bind.KeyServer(cfg, cfg.KeyEndpoint())
	store, _ = bind.StoreServer(cfg, cfg.KeyEndpoint())
	dir = inprocess.New(cfg)
	return
}

func makeDirectory(dir upspin.DirServer, directoryName upspin.PathName) (*upspin.DirEntry, error) {
	parsed, err := path.Parse(directoryName)
	if err != nil {
		return nil, err
	}
	// Can't use newDirEntry as it adds a block.
	entry := &upspin.DirEntry{
		Name:       parsed.Path(),
		SignedName: parsed.Path(),
		Attr:       upspin.AttrDirectory,
	}
	return dir.Put(entry)
}

func Test_Lookup(t *testing.T) {
	cfg, _, _, _ := newConfigAndServices("test@gmail.com")
	ipdir := inprocess.New(cfg)
	_, err := makeDirectory(ipdir, upspin.PathName("test@gmail.com"))
	if err != nil {
		panic(err)
	}

	_, err = ipdir.Put(&upspin.DirEntry{
		Name:       upspin.PathName("test@gmail.com/toto"),
		SignedName: upspin.PathName("test@gmail.com/toto"),
		Packing:    upspin.PlainPack,
		Writer:     upspin.UserName("test@gmail.com"),
	})
	fmt.Println(err)
	assert.NoError(t, err)

	expected, err := ipdir.Lookup("test@gmail.com/toto")
	assert.NoError(t, err)

	actual, err := (&Dir{
		Username: "test@gmail.com",
		Root:     "test_dir",
	}).Lookup("test@gmail.com/toto")
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

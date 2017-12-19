package packing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"upspin.io/bind"
	"upspin.io/config"
	"upspin.io/factotum"
	"upspin.io/pack"
	_ "upspin.io/pack/plain"
	"upspin.io/test/testutil"
	_ "upspin.io/transports"
	"upspin.io/upspin"
)

func Test_PlainPackRecognizedByUnpack(t *testing.T) {
	cfg, _, _, _ := newConfigAndServices(upspin.UserName("gildaschbt+local@gmail.com"))

	f, _ := os.Open("albert.txt")
	fi, _ := f.Stat()

	d := PlainDirEntry(fi, cfg)

	_, err := pack.Lookup(upspin.PlainPack).Unpack(cfg, d)

	assert.NoError(t, err)
}

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
	dir, _ = bind.DirServer(cfg, cfg.KeyEndpoint())
	return
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gildasch/upspin-localserver/dir"
	"github.com/gildasch/upspin-localserver/local"
	"github.com/gildasch/upspin-localserver/store"
	"upspin.io/config"
	"upspin.io/factotum"
	_ "upspin.io/key/transports"
	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	"upspin.io/upspin"
)

func main() {
	rootPtr := flag.String("root", ".",
		"the root directory to serve")
	debugPtr := flag.Bool("debug", false,
		"activate debug mode")
	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := upspin.NetAddr("http://localhost:" + port)

	cfg := newConfig()

	dirServer := dirserver.New(
		cfg,
		&dir.Dir{
			Username: "gildaschbt+local@gmail.com",
			Root:     *rootPtr,
			Storage:  &local.Storage{*rootPtr},
			Debug:    *debugPtr,
			Factotum: cfg.Factotum()},
		addr)

	http.Handle("/api/Dir/", dirServer)

	storeServer := storeserver.New(
		cfg,
		&store.Store{
			Root:  *rootPtr,
			Debug: *debugPtr},
		addr)

	http.Handle("/api/Store/", storeServer)

	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func newConfig() upspin.Config {
	endpoint := upspin.Endpoint{
		Transport: upspin.Remote,
		NetAddr:   "usl.gildas.ch",
	}
	cfg := config.New()
	cfg = config.SetUserName(cfg, upspin.UserName("gildaschbt+local@gmail.com"))
	cfg = config.SetPacking(cfg, upspin.PlainPack)
	cfg = config.SetStoreEndpoint(cfg, endpoint)
	cfg = config.SetDirEndpoint(cfg, endpoint)

	f, err := factotum.NewFromDir(
		"/home/gildas/.ssh/gildaschbt+local@gmail.com")
	if err != nil {
		panic(err)
	}
	cfg = config.SetFactotum(cfg, f)

	return cfg
}

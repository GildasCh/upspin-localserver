package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gildasch/upspin-localserver/dir"
	"upspin.io/config"
	_ "upspin.io/key/transports"
	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	"upspin.io/upspin"
)

type store struct {
	upspin.StoreServer
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := upspin.NetAddr("http://localhost:" + port)

	dirCfg := config.New()
	ep := upspin.Endpoint{
		Transport: upspin.Remote,
		NetAddr:   "usl.gildas.ch",
	}
	dirCfg = config.SetDirEndpoint(dirCfg, ep)
	dirCfg = config.SetStoreEndpoint(dirCfg, ep)

	dirServer := dirserver.New(
		dirCfg,
		&dir.Dir{},
		addr)

	http.Handle("/api/Dir/", dirServer)

	storeServer := storeserver.New(
		dirCfg,
		store{},
		addr)

	http.Handle("/api/Store/", storeServer)

	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

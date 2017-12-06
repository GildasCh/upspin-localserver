package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gildasch/upspin-localserver/dir"
	"github.com/gildasch/upspin-localserver/store"
	"upspin.io/config"
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

	dirCfg := config.New()
	ep := upspin.Endpoint{
		Transport: upspin.Remote,
		NetAddr:   "usl.gildas.ch",
	}
	dirCfg = config.SetDirEndpoint(dirCfg, ep)
	dirCfg = config.SetStoreEndpoint(dirCfg, ep)

	dirServer := dirserver.New(
		dirCfg,
		&dir.Dir{
			Username: "gildaschbt+local@gmail.com",
			Root:     *rootPtr,
			Debug:    *debugPtr},
		addr)

	http.Handle("/api/Dir/", dirServer)

	storeServer := storeserver.New(
		dirCfg,
		&store.Store{
			Root:  *rootPtr,
			Debug: *debugPtr},
		addr)

	http.Handle("/api/Store/", storeServer)

	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

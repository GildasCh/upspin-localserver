package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gildasch/upspin-localserver/dir"
	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	"upspin.io/upspin"
)

type config struct {
	upspin.Config
}

type store struct {
	upspin.StoreServer
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := upspin.NetAddr("http://localhost:" + port)

	dirServer := dirserver.New(
		config{},
		&dir.Dir{},
		addr)

	http.Handle("/dir", dirServer)

	storeServer := storeserver.New(
		config{},
		store{},
		addr)

	http.Handle("/store", storeServer)

	fmt.Printf("Listening on %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

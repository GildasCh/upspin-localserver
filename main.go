package main

import (
	"fmt"
	"net/http"
	"os"

	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	"upspin.io/upspin"
)

type config struct {
	upspin.Config
}

type dir struct {
	upspin.DirServer
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
		dir{},
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

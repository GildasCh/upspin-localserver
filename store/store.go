package store

import (
	"fmt"

	"upspin.io/upspin"
)

type Store struct {
	upspin.StoreServer

	Root string
}

func (s *Store) Dial(upspin.Config, upspin.Endpoint) (upspin.Service, error) {
	return s, nil
}

func (s *Store) Endpoint() upspin.Endpoint {
	return upspin.Endpoint{}
}

func (s *Store) Close() {}

func (s *Store) Get(ref upspin.Reference) ([]byte, *upspin.Refdata, []upspin.Location, error) {
	fmt.Println("requested ref:", ref)
	return []byte("hello"), nil, nil, nil
}

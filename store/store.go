package store

import (
	"fmt"

	"upspin.io/upspin"
)

type Store struct {
	upspin.StoreServer

	Root  string
	Debug bool
}

func (s *Store) Dial(config upspin.Config, endpoint upspin.Endpoint) (upspin.Service, error) {
	if s.Debug {
		fmt.Printf("dir.Dial called with config=%#v, endpoint=%#v\n", config, endpoint)
	}

	return s, nil
}

func (s *Store) Endpoint() upspin.Endpoint {
	if s.Debug {
		fmt.Printf("store.Endpoint called\n")
	}

	return upspin.Endpoint{}
}

func (s *Store) Close() {
	if s.Debug {
		fmt.Printf("store.Close called\n")
	}
}

func (s *Store) Get(ref upspin.Reference) ([]byte, *upspin.Refdata, []upspin.Location, error) {
	if s.Debug {
		fmt.Printf("store.Get called with ref=%#v\n", ref)
	}

	fmt.Println("requested ref:", ref)

	if s.Debug {
		fmt.Printf("store.Get returning %#v\n", []byte("hello"))
	}

	return []byte("hello"),
		&upspin.Refdata{Reference: ref},
		nil,
		nil
}

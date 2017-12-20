package store

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"upspin.io/errors"
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

	if ref == upspin.HTTPBaseMetadata {
		return nil, nil, nil, errors.E(errors.NotExist)
	}

	split := strings.Split(string(ref), "-")

	offset, err := strconv.ParseInt(split[len(split)-1], 10, 64)
	if err != nil {
		return nil, nil, nil, errors.E(errors.NotExist)
	}

	relativePath := strings.Join(
		split[:len(split)-1], "")
	f, err := os.Open(path.Join(s.Root, relativePath))
	if err != nil {
		return nil, nil, nil, errors.E(errors.NotExist)
	}

	bytes := make([]byte, upspin.BlockSize)
	_, err = f.ReadAt(bytes, offset)
	if err != nil && err != io.EOF {
		return nil, nil, nil, errors.E(errors.IO)
	}

	if s.Debug {
		fmt.Printf("store.Get returning byte array of lenght %d starting with %#v\n", len(bytes), bytes[:20])
	}

	return bytes, &upspin.Refdata{Reference: ref}, nil, nil
}

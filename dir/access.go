package dir

import (
	"errors"
	"fmt"

	"upspin.io/access"
	"upspin.io/path"
	"upspin.io/upspin"
)

func (d *Dir) canList(dir upspin.PathName) bool {
	if d.userName == "" {
		return false
	}

	accessPath := path.Clean(dir + "/Access")
	if accData, ok := d.Storage.Access(string(accessPath)); ok {
		acc, err := access.Parse(accessPath, accData)
		if err != nil {
			fmt.Println(err)
			return false
		}
		can, err := acc.Can(d.userName, access.List, dir, func(upspin.PathName) ([]byte, error) { return nil, errors.New("") })
		if err != nil {
			fmt.Println(err)
			return false
		}
		return can
	}

	dirParsed, err := path.Parse(dir)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if dirParsed.IsRoot() {
		return false
	}
	return d.canList(dirParsed.Drop(1).Path())
}

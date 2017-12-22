package packing

import (
	"github.com/gildasch/upspin-localserver/local"
	"upspin.io/upspin"
)

type Simulator interface {
	DirEntry(username string, fi local.FileInfo, factotum Factotum) *upspin.DirEntry
}

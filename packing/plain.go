package packing

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/gildasch/upspin-localserver/local"
	"upspin.io/pack/packutil"
	"upspin.io/upspin"
)

const (
	aesKeyLen     = 32
	marshalBufLen = 66
)

var (
	zero = big.NewInt(0)
)

type Factotum interface {
	FileSign(hash upspin.DEHash) (upspin.Signature, error)
	DirEntryHash(
		n, l upspin.PathName, a upspin.Attribute, p upspin.Packing,
		t upspin.Time, dkey, hash []byte) upspin.DEHash
}

type Plain struct{}

func (Plain) DirEntry(username string, fi local.FileInfo, factotum Factotum) *upspin.DirEntry {
	e := dirEntryFromFileInfo(username, fi)

	// Compute entry signature with dkey=sum=0.
	dkey := make([]byte, aesKeyLen)
	sum := make([]byte, sha256.Size)
	sig, err := factotum.FileSign(factotum.DirEntryHash(e.SignedName, e.Link, e.Attr, e.Packing, e.Time, dkey, sum))
	if err != nil {
		panic(err.Error())
	}

	pdMarshal(&e.Packdata, sig, upspin.Signature{})

	return e
}

func dirEntryFromFileInfo(username string, fi local.FileInfo) *upspin.DirEntry {
	de := &upspin.DirEntry{
		Name: upspin.PathName(
			username + fi.Filename),
		Packing: upspin.PlainPack,
		Writer:  upspin.UserName(username),
		Time:    upspin.Time(fi.Time.Unix()),
	}
	if fi.IsDir {
		de.Attr = upspin.AttrDirectory
	} else {
		de.Blocks = blocksFromFileInfo(fi)
	}
	return de
}

func blocksFromFileInfo(fi local.FileInfo) (dbs []upspin.DirBlock) {
	size := fi.Size
	offset := int64(0)
	for size > 0 {
		s := int64(upspin.BlockSize)
		if s > size {
			s = size
		}
		size -= s
		ref := fmt.Sprintf("%s-%d", fi.Filename, offset)
		dbs = append(dbs, upspin.DirBlock{
			Location: upspin.Location{
				Endpoint: upspin.Endpoint{
					Transport: upspin.Remote,
					NetAddr:   "usl.gildas.ch",
				},
				Reference: upspin.Reference(ref)},
			Offset: offset,
			Size:   s,
		})
		offset += s
	}

	return
}

func pdMarshal(dst *[]byte, sig, sig2 upspin.Signature) error {
	// sig2 is a signature with another owner key, to enable smoother key rotation.
	n := packdataLen()
	if len(*dst) < n {
		*dst = make([]byte, n)
	}
	n = 0
	n += packutil.PutBytes((*dst)[n:], sig.R.Bytes())
	n += packutil.PutBytes((*dst)[n:], sig.S.Bytes())
	if sig2.R == nil {
		sig2 = upspin.Signature{R: zero, S: zero}
	}
	n += packutil.PutBytes((*dst)[n:], sig2.R.Bytes())
	n += packutil.PutBytes((*dst)[n:], sig2.S.Bytes())
	*dst = (*dst)[:n]
	return nil
}

// packdataLen returns n big enough for packing, sig.R, sig.S
func packdataLen() int {
	return 2*marshalBufLen + binary.MaxVarintLen64 + 1
}

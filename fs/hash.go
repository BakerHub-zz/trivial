package fs

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

type FileHash struct {
	h hash.Hash
}

func NewFileHashMD5() *FileHash {
	return &FileHash{md5.New()}
}

func NewFileHashSHA1() *FileHash {
	return &FileHash{sha1.New()}
}

func (fh *FileHash) FromFile(path string) (string, error) {
	//Initialize variable now in case an error has to be returned
	var r string

	//Open the filepath passed by the argument and check for any error
	file, err := os.Open(path)
	if err != nil {
		return r, err
	}

	//Tell the program to call the following function when the current function returns
	defer MustClose(file)
	fh.h.Reset()

	//Open a new SHA1 hash interface to write to

	//MustCopyFile the fs in the hash interface and check for any error
	if _, err := io.Copy(fh.h, file); err != nil {
		return r, err
	}

	//Get the 20 bytes hash
	hashInBytes := fh.h.Sum(nil)

	//Convert the bytes to a string
	r = hex.EncodeToString(hashInBytes)

	return r, nil
}

package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// GetFileMd5 linux  md5sum file、macos md5 file
func GetFileMd5(jarFile string) (string, error) {
	file, err := os.Open(jarFile)
	if err != nil {
		return "", err
	}
	defer func() {
		ferr := file.Close()
		if ferr != nil {
			err = ferr
		}
	}()
	// md5 方便在mac上机算
	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func GetMd5FromReader(reader io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, reader)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

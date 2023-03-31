package util

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
)

type Manifest map[string]string

var ErrNotJAR = errors.New("Given file is not a JAR file")
var ErrWrongManifestFormat = errors.New("Can't parse manifest file (wrong format)")

// ReadFile reads JAR file and parses manifest file
func ReadFile(jarFile string) (Manifest, error) {
	r, err := zip.OpenReader(jarFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name != "META-INF/MANIFEST.MF" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		return readManifestData(rc)
	}
	return nil, ErrNotJAR
}

func readManifestData(r io.Reader) (Manifest, error) {
	m := make(Manifest)
	s := bufio.NewScanner(r)
	var propName, propVal string
	for s.Scan() {
		text := s.Text()
		if len(text) == 0 {
			continue
		}
		if strings.HasPrefix(text, " ") {
			m[propName] += strings.TrimLeft(text, " ")
			continue
		}
		propSepIndex := strings.Index(text, ": ")
		if propSepIndex == -1 || len(text) < propSepIndex+2 {
			return nil, ErrWrongManifestFormat
		}
		propName = text[:propSepIndex]
		propVal = text[propSepIndex+2:]
		m[propName] = propVal
	}
	return m, nil
}

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

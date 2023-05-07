package common

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"server/resources"
	"strings"
)

type Resource struct {
	fs   embed.FS
	path string
}

func NewResource() *Resource {
	return &Resource{
		fs:   resources.Static,
		path: "html",
	}
}

func (r *Resource) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := filepath.Join(r.path, filepath.FromSlash(path.Clean("/static/"+name)))
	fullName = filepath.ToSlash(fullName)
	file, err := r.fs.Open(fullName)

	return file, err
}

package util

import (
	"archive/zip"
	"io"
	"server/model"
)

func GetZipItemInfo(zipFile string) ([]model.ZipItemInfo, error) {
	var result = []model.ZipItemInfo{}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return result, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		reader, err := f.Open()
		if err != nil {
			continue
		}
		md5, err := GetMd5FromReader(reader)
		reader.Close()
		item := model.ZipItemInfo{
			FileName: f.Name,
			Md5:      md5,
		}
		result = append(result, item)
	}
	return result, nil
}

func ReadFileFromZipByPath(zipName string, filePath string) ([]byte, error) {
	r, err := zip.OpenReader(zipName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	reader, err := r.Open(filePath)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

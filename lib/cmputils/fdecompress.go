package cmputil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

//UnGzip decompresses a file compressed using gzip
func (f *File) UnGzip(target string) error {
	reader, err := os.Open(f.Path)

	if err != nil {
		return err
	}

	defer reader.Close()

	archive, err := gzip.NewReader(reader)

	if err != nil {
		return err
	}

	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)

	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

//Untar decompresses a file compressed using tar
func (f *File) Untar(target string) error {
	reader, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

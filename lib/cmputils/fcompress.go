package cmputil

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//File represents a file at path: Path
type File struct {
	Path string
}

//Gzip compresses a file using gzip compression, and places the result at [target]
func (f *File) Gzip(target string) error {
	reader, opErr := os.Open(f.Path)

	if opErr != nil {
		fmt.Println(opErr)
		return opErr
	}

	filename := filepath.Base(f.Path)

	if !(filepath.Ext(f.Path) == "") {
		filename = strings.Replace(filename, filepath.Ext(f.Path), "", 1)
	}

	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)

	if err != nil {
		return err
	}

	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename

	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

//Tar compresses a file using tar compression, and places the result at [target]
func (f *File) Tar(target string) error {
	filename := filepath.Base(f.Path)

	if !(filepath.Ext(f.Path) == "") {
		filename = strings.Replace(filename, filepath.Ext(f.Path), "", 1)
	}

	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(f.Path)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(f.Path)
	}

	return filepath.Walk(f.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, f.Path))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

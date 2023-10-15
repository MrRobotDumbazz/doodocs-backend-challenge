package zipfile

import (
	"archive/zip"
	"fmt"
)

type ZipFile struct {
	zipfile *zip.ReadCloser
}

func New(filename string) (*zip.ReadCloser, error) {
	const op = "zipfile.NewZipFile"
	zipfile, err := zip.OpenReader(filename)
	defer zipfile.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return zipfile, nil
}

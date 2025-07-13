package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type FS struct {
	Base string
}

func (f FS) Save(ctx context.Context, name string, r io.Reader) (string, error) {
	path := filepath.Join(f.Base, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, r); err != nil {
		return "", err
	}
	return "/uploads/" + name, nil
}

func (f FS) URL(name string) string {
	return "/uploads/" + name
}

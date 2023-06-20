package src

import (
	"errors"
	"os"
	"path/filepath"
)

func New(srcPath string) (Sourcer, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(srcPath)

	switch ext {
	case ".csv":
		return &csvSourcer{file: f}, nil
	case ".json":
		return &jsonSourcer{file: f}, nil
	default:
		return nil, errors.New("unsupported source file type")
	}
}

type Sourcer interface {
	Data() (map[string]string, error)
}

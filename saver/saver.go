package saver

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Saver interface {
	Save(map[string]string) error
}

func New(path, ext string) (Saver, error) {

	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	switch ext {
	case ".csv":
		return &csvSaver{path: path}, nil
	case ".json":
		return &jsonSaver{path: path}, nil
	default:
		return nil, errors.New("unsupported file type for saving")

	}

}

func createFile(path, extension string) (*os.File, error) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	dirPath := filepath.Dir(path)

	err := os.MkdirAll(dirPath, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	return os.Create(fileName(dirPath, extension))
}

func fileName(dirPath, extension string) string {

	fname := time.Now().Format("save_2006-01-02-15")

	return fmt.Sprintf("%s/%s.%s", dirPath, fname, extension)
}

package saver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Saver interface {
	Save(map[string]string) error
}

type CreateSaverFunc func(output string) Saver

var Factory = map[string]CreateSaverFunc{
	"csv":  newCsvSaver,
	"json": newJsonSaver,
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

package saver

import (
	"encoding/json"
)

func newJsonSaver(path string) Saver {
	return &jsonSaver{path: path}
}

type jsonSaver struct {
	path string
}

func (s *jsonSaver) Save(records map[string]string) error {
	f, err := createFile(s.path, "json")
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(records)
}

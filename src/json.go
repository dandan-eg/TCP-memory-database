package src

import (
	"encoding/json"
	"os"
)

type jsonSourcer struct {
	file *os.File
}

func (j *jsonSourcer) Data() (map[string]string, error) {
	defer j.file.Close()

	data := make(map[string]string)

	if err := json.NewDecoder(j.file).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil

}

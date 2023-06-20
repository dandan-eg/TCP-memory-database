package src

import (
	"encoding/csv"
	"errors"
	"os"
)

type csvSourcer struct {
	file *os.File
}

func (c csvSourcer) Data() (map[string]string, error) {
	defer c.file.Close()
	data := make(map[string]string)

	r := csv.NewReader(c.file)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) != 2 {
			return nil, errors.New("record must be key/value paired")
		}

		k, v := record[0], record[1]
		data[k] = v
	}

	return data, nil
}

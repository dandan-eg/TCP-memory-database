package saver

import (
	"encoding/csv"
)

type csvSaver struct {
	path string
}

func newCsvSaver(path string) Saver {
	return &csvSaver{path}
}

func (s *csvSaver) Save(records map[string]string) error {
	f, err := createFile(s.path, "csv")
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for k, v := range records {
		ln := []string{k, v}

		err := w.Write(ln)
		if err != nil {
			return err
		}
	}

	return nil
}

package codecs

import (
	"encoding/csv"
	"io"
)

func CsvReaderFunc(r io.Reader) ([][]string, error) {
	return csv.NewReader(r).ReadAll()
}

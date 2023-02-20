package postgresql

import (
	"strconv"
	"strings"
)

func ToArray(headers []string, rows [][]string) string {
	columnsCount := len(rows[0])
	transposed := make([][]string, columnsCount)
	for _, row := range rows {
		for n, column := range row {
			transposed[n] = append(transposed[n], "'"+column+"'")
		}
	}

	var arrays []string
	for _, arrayList := range transposed {
		arrays = append(arrays, "ARRAY ["+strings.Join(arrayList, ",")+"]")
	}

	var headerColumns []string
	if len(headers) == 0 {
		for i := 1; i <= columnsCount; i++ {
			headerColumns = append(headerColumns, "\"col"+strconv.Itoa(i)+"\"")
		}
	} else {
		for _, header := range headers {
			headerColumns = append(headerColumns, "\""+header+"\"")
		}
	}

	return "unnest(" + strings.Join(arrays, ",\n") + ")" +
		"\nAS csv(" + strings.Join(headerColumns, ",") + ")"
}

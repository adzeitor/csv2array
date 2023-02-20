package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

var (
	ErrHeaderExpected = errors.New("header expected")
)

type Config struct {
	SkipHeader bool
}

type Magi struct {
	Config   Config
	Template Converter
}

func New(tmpl Converter, cfg Config) (Magi, error) {
	magi := Magi{
		Config:   cfg,
		Template: tmpl,
	}
	return magi, nil
}

type header = []string
type rows = [][]string
type Converter = func(header, rows) string

func (magi *Magi) ReadAndExecute(r io.Reader, w io.Writer) error {
	csvReader := csv.NewReader(r)

	var headers []string
	// TODO: add header to converter
	if !magi.Config.SkipHeader {
		h, err := readHeader(csvReader)
		if err != nil {
			return err
		}
		headers = h
	}

	var rows [][]string
	for line := 1; ; line++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("line %d: error %w\n", line, err)
		}

		rows = append(rows, record)
	}
	result := magi.Template(headers, rows)

	_, err := w.Write([]byte(result + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func readHeader(r *csv.Reader) (header, error) {
	record, err := r.Read()
	if err == io.EOF {
		return nil, ErrHeaderExpected
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

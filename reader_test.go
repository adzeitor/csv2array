package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/adzeitor/csv2array/dialects/postgresql"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		in      string
		out     string
		wantErr bool
	}{
		{
			name: "success",
			in: "col1,col 2,col+3\n" +
				"one,two,three\n" +
				"1,2,3\n",
			out: "unnest(" +
				"ARRAY ['one','1'],\n" +
				"ARRAY ['two','2'],\n" +
				"ARRAY ['three','3'])\n" +
				"AS csv(\"col1\",\"col 2\",\"col+3\")\n",
		},
		{
			name: "without header",
			cfg:  Config{SkipHeader: true},
			in: "one,two,three\n" +
				"1,2,3\n",
			out: "unnest(" +
				"ARRAY ['one','1'],\n" +
				"ARRAY ['two','2'],\n" +
				"ARRAY ['three','3'])\n" +
				"AS csv(\"col1\",\"col2\",\"col3\")\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			magi, err := New(postgresql.ToArray, tt.cfg)
			if err != nil {
				t.Fatal(err)
			}

			err = magi.ReadAndExecute(strings.NewReader(tt.in), buf)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Expected error but got %v\n", err)
			}

			got := buf.String()
			if got != tt.out {
				t.Fatalf("expected %q, but got %q", tt.out, got)
			}
		})
	}
}

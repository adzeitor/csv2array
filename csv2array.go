package main

import (
	"flag"
	"log"
	"os"

	"github.com/adzeitor/csv2array/dialects/postgresql"
)

func main() {
	in := os.Stdin
	out := os.Stdout
	skipHeader := flag.Bool("skipheader", false, "use strict mode (like error on undefined columns)")
	dialect := flag.String("dialect", "postgresql", "dialect")
	flag.Parse()

	var converter Converter
	switch *dialect {
	case "postgresql":
		converter = postgresql.ToArray
	default:
		panic("unknown dialect")
	}

	magi, err := New(converter, Config{
		SkipHeader: *skipHeader,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = magi.ReadAndExecute(in, out)
	if err != nil {
		log.Fatalln(err)
	}
}

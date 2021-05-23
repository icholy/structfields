package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/icholy/structfields"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", "working directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [packages]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	ss, err := structfields.Load(dir, flag.Args()...)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(ss); err != nil {
		log.Fatal(err)
	}
}

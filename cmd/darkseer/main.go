package main

import (
	"flag"
	"log"
	"os"

	"github.com/aoisensi/darkseer/dmx"
	"github.com/k0kubun/pp/v3"
)

func main() {
	flag.Parse()
	for _, name := range flag.Args() {
		f, err := os.Open(name)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		element, err := dmx.NewDecoder(f).Decode()
		if err != nil {
			log.Println(err)
		}
		pp.Println(element)
	}
}

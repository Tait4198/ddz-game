package main

import (
	"flag"
)

func main() {
	var name string
	flag.StringVar(&name, "n", "cyt", "")
	flag.Parse()

	dc := NewDdzClient(name, "123")
	dc.Run()
}

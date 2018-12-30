package main

import (
	"flag"
	"simon/iface"
)

func main() {
	flag.Parse()
	iface.CLI.Start()
}

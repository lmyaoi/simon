package main

import (
	"flag"
	"vsync/iface"
)

func main() {
	flag.Parse()
	iface.CLI.Start()
}

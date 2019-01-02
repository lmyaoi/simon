package main

import (
	"simon/iface"
	"simon/conf"
)

func init() {
	conf.Initialize()
}

func main() {
	iface.CLI.Start()
}

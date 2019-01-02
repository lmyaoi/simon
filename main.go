package main

import (
	"simon/iface"
	"simon/pref"
)

func init() {
	pref.Initialize()
}

func main() {
	iface.CLI.Start()
}

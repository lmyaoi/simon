//go:generate goversioninfo

package main

import (
	"simon/conf"
	"simon/iface"
)

func init() {
	conf.Initialize()
}

func main() {
	iface.CLI.Start()
}

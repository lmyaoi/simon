package iface

import "vsync/iface/cli"

type Interface interface {
	Start()
}

const CLI cli.CLI = 0
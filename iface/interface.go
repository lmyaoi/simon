package iface

import "simon/iface/cli"

type Interface interface {
	Start()
}

const CLI cli.CLI = 0

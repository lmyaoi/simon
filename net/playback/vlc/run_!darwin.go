// +build !darwin

package vlc

import (
	"os/exec"
	"simon/conf"
	"strconv"
)

func run() {
	port := strconv.Itoa(conf.Get().VlcPort())
	args := []string{"--extraintf", "http", "--http-port", port, "--http-password", "q"}
	cmd := exec.Command(conf.Get().VlcPath(), args...)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

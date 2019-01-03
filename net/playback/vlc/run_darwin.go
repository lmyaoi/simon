package vlc

import (
	"os/exec"
	"simon/conf"
	"strconv"
)

func run() {
	port := strconv.Itoa(conf.Get().VlcPort)
	args := []string{"-a", "vlc", "--args", "--extraintf", "http", "--http-port", port, "--http-password", "q"}
	cmd := exec.Command("open", args...)
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

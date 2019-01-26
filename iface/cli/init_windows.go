package cli

import (
	"golang.org/x/sys/windows"
	"os"
)

func init() {
	handle := windows.Handle(os.Stdout.Fd())
	var originalMode uint32
	err := windows.GetConsoleMode(handle, &originalMode)
	if err != nil {
		return
	}
	err = windows.SetConsoleMode(handle, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if err != nil {
		panic(err)
	}
}

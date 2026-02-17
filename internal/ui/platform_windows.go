//go:build windows

// +build windows

package ui

import (
	"golang.org/x/sys/windows"
	"os"
)

func EnableVirtualTerminalProcessing() {
	handle := windows.Handle(os.Stdout.Fd())
	var mode uint32
	if err := windows.GetConsoleMode(handle, &mode); err == nil {
		mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		windows.SetConsoleMode(handle, mode)
	}
}

//go:build windows

package graceful

import (
	"fmt"
	"syscall"
)

// sendSignal sends a signal to a given process id.
func sendSignal(pid int, _ syscall.Signal) error {
	d, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return fmt.Errorf("failed to load kernel32.dll: %w", err)
	}
	p, err := d.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return fmt.Errorf("failed to find GenerateConsoleCtrlEvent: %w", err)
	}
	r, _, err := p.Call(uintptr(syscall.CTRL_CLOSE_EVENT), uintptr(pid))
	if r == 0 {
		return fmt.Errorf("failed to generate console control event: %w", err)
	}
	return nil
}

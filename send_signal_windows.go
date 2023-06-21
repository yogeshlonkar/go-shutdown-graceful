//go:build windows

package graceful

import (
	"fmt"
	"os"
	"syscall"
)

// sendSignal sends a signal to a given process id.
func sendSignal(pid int, sig syscall.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}
	if err = process.Signal(sig); err != nil {
		return fmt.Errorf("failed to send signal %s(%#v) to current process %d: %w", sig.String(), sig, pid, err)
	}
	return nil
}

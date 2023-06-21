//go:build !windows

package graceful

import (
	"fmt"
	"syscall"
)

// sendSignal will send a signal to given process id.
func sendSignal(pid int, sig syscall.Signal) error {
	err := syscall.Kill(pid, sig)
	if err != nil {
		return fmt.Errorf("failed to send signal %s(%#v) to current process %d: %w", sig.String(), sig, pid, err)
	}
	return nil
}

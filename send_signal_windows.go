//go:build windows

package graceful

import (
	"errors"
	"fmt"
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

// sendSignal sends a signal to a given process id.
func sendSignal(pid int, _ syscall.Signal) error {
	// process, err := os.FindProcess(pid)
	// if err != nil {
	//	return fmt.Errorf("failed to find process %d: %w", pid, err)
	//}
	// if err = process.Signal(sig); err != nil {
	//	return fmt.Errorf("failed to send signal %s(%#v) to current process %d: %w", sig.String(), sig, pid, err)
	//}
	// return nil
	dll, err := windows.LoadDLL("kernel32.dll")
	if err != nil {
		return fmt.Errorf("failed to load kernel32.dll: %w", err)
	}
	defer func(dll *windows.DLL) {
		if err = dll.Release(); err != nil {
			log.Printf("send CTRL_CLOSE_EVENT: failed to release kernel32.dll: %v", err)
		}
	}(dll)
	f, err := dll.FindProc("AttachConsole")
	if err != nil {
		return fmt.Errorf("failed to find AttachConsole: %w", err)
	}
	r1, _, err := f.Call(uintptr(pid))
	if r1 == 0 && errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
		return fmt.Errorf("failed to attach console: %w", err)
	}
	f, err = dll.FindProc("SetConsoleCtrlHandler")
	if err != nil {
		return fmt.Errorf("failed to find SetConsoleCtrlHandler: %w", err)
	}
	r1, _, err = f.Call(0, 1)
	if r1 == 0 {
		return fmt.Errorf("failed to set console control handler: %w", err)
	}
	f, err = dll.FindProc("GenerateConsoleCtrlEvent")
	if err != nil {
		return fmt.Errorf("failed to find GenerateConsoleCtrlEvent: %w", err)
	}
	r1, _, err = f.Call(windows.CTRL_CLOSE_EVENT, uintptr(pid))
	if r1 == 0 {
		return fmt.Errorf("failed to generate console control event: %w", err)
	}
	return nil
}

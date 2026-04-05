package killer

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// KillResult holds the outcome of a kill attempt.
type KillResult struct {
	PID     int
	Name    string
	Port    int
	Success bool
	Err     error
}

// KillPID sends SIGTERM to pid, waits up to 2s, then SIGKILL if still alive.
// Returns nil on success, error on failure.
func KillPID(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("process %d not found: %w", pid, err)
	}

	// Send SIGTERM
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		if isPermission(err) {
			return fmt.Errorf("permission denied — try running with sudo")
		}
		return fmt.Errorf("failed to send SIGTERM to PID %d: %w", pid, err)
	}

	// Poll up to 2s for the process to exit
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(100 * time.Millisecond)
		if !isAlive(pid) {
			return nil
		}
	}

	// Still alive — send SIGKILL
	if err := proc.Signal(syscall.SIGKILL); err != nil {
		if isPermission(err) {
			return fmt.Errorf("permission denied — try running with sudo")
		}
		return fmt.Errorf("failed to kill PID %d: %w", pid, err)
	}
	// SIGKILL is unblockable; if it sent without error, the process is done.
	return nil
}

// isAlive checks if a process with the given PID is still running.
func isAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}

func isPermission(err error) bool {
	return err == syscall.EPERM
}

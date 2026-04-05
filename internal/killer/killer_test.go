package killer_test

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/mohokh67/portly/internal/killer"
)

func TestKillPID_AlreadyDead(t *testing.T) {
	err := killer.KillPID(99999999)
	if err == nil {
		t.Error("expected error killing nonexistent PID")
	}
}

func TestKillPID_LiveProcess(t *testing.T) {
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start test process: %v", err)
	}
	pid := cmd.Process.Pid

	err := killer.KillPID(pid)
	if err != nil {
		t.Fatalf("KillPID failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	proc, err := os.FindProcess(pid)
	if err == nil {
		_ = proc.Signal(os.Signal(nil))
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case <-done:
		// good
	case <-time.After(3 * time.Second):
		t.Error("process still alive after KillPID")
	}
}

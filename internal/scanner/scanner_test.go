package scanner_test

import (
	"testing"

	"github.com/mohokh67/portly/internal/scanner"
)

func TestScanReturnsProcesses(t *testing.T) {
	procs, err := scanner.Scan(scanner.ListeningOnly)
	if err != nil {
		t.Fatalf("Scan() error: %v", err)
	}
	// On any real machine there should be at least one listening port
	if len(procs) == 0 {
		t.Skip("no listening ports found — skipping (CI may have none)")
	}
	p := procs[0]
	if p.Port <= 0 || p.Port > 65535 {
		t.Errorf("invalid port: %d", p.Port)
	}
	if p.PID <= 0 {
		t.Errorf("invalid PID: %d", p.PID)
	}
	if p.Proto != "TCP" && p.Proto != "UDP" {
		t.Errorf("unexpected proto: %s", p.Proto)
	}
}

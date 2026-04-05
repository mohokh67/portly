package scanner

import (
	"testing"
)

func TestParseLsofOutput(t *testing.T) {
	sample := `COMMAND     PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
node      98234   mehr   23u  IPv4 0x1234      0t0  TCP *:3000 (LISTEN)
postgres   1204 _postgres  5u  IPv4 0x5678      0t0  TCP 127.0.0.1:5432 (LISTEN)
node      98234   mehr   24u  IPv4 0x9abc      0t0  TCP 127.0.0.1:3000->127.0.0.1:54321 (ESTABLISHED)
`
	procs, err := parseLsof(sample, ListeningOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(procs) != 2 {
		t.Fatalf("expected 2 listening procs, got %d", len(procs))
	}

	if procs[0].Port != 3000 {
		t.Errorf("expected port 3000, got %d", procs[0].Port)
	}
	if procs[0].Name != "node" {
		t.Errorf("expected name 'node', got %q", procs[0].Name)
	}
	if procs[0].PID != 98234 {
		t.Errorf("expected PID 98234, got %d", procs[0].PID)
	}
	if procs[0].Address != "0.0.0.0" {
		t.Errorf("expected addr 0.0.0.0, got %q", procs[0].Address)
	}
	if procs[1].Address != "127.0.0.1" {
		t.Errorf("expected addr 127.0.0.1, got %q", procs[1].Address)
	}
}

func TestParseLsofAllConnections(t *testing.T) {
	sample := `COMMAND     PID   USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
node      98234   mehr   23u  IPv4 0x1234      0t0  TCP *:3000 (LISTEN)
node      98234   mehr   24u  IPv4 0x9abc      0t0  TCP 127.0.0.1:3000->127.0.0.1:54321 (ESTABLISHED)
`
	procs, err := parseLsof(sample, AllConnections)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(procs) != 2 {
		t.Fatalf("expected 2 procs in AllConnections mode, got %d", len(procs))
	}
}

func TestParseLsofEmptyOutput(t *testing.T) {
	procs, err := parseLsof("", ListeningOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(procs) != 0 {
		t.Errorf("expected 0 procs, got %d", len(procs))
	}
}

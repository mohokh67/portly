package scanner

import (
	"testing"
)

func TestParseProcNetTCP(t *testing.T) {
	// /proc/net/tcp format: sl local_address rem_address st ...
	// 0100007F:0CEA = 127.0.0.1:3306  (0CEA hex = 3306)
	// 00000000:0BB8 = 0.0.0.0:3000    (0BB8 hex = 3000)
	// state 0A = LISTEN
	sample := `  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
   0: 00000000:0BB8 00000000:0000 0A 00000000:00000000 00:00000000 00000000  1000        0 12345 1 0000000000000000 100 0 0 10 0
   1: 0100007F:0CEA 00000000:0000 0A 00000000:00000000 00:00000000 00000000   26        0 67890 1 0000000000000000 100 0 0 10 0
`
	entries, err := parseProcNetTCP(sample, ListeningOnly)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 3000 {
		t.Errorf("expected port 3000, got %d", entries[0].Port)
	}
	if entries[0].Address != "0.0.0.0" {
		t.Errorf("expected 0.0.0.0, got %q", entries[0].Address)
	}
	if entries[1].Port != 3306 {
		t.Errorf("expected port 3306, got %d", entries[1].Port)
	}
	if entries[1].Address != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1, got %q", entries[1].Address)
	}
}

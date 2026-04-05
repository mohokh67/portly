package scanner

import "errors"

// ScanMode controls which connections are returned.
type ScanMode int

const (
	ListeningOnly  ScanMode = iota // only LISTEN state (default)
	AllConnections                 // includes ESTABLISHED, etc.
)

// Process represents a single process bound to a port.
type Process struct {
	Port    int
	Proto   string // "TCP" or "UDP"
	PID     int
	User    string
	Address string // e.g. "0.0.0.0" or "127.0.0.1"
	Name    string // process name (e.g. "node")
}

// Scan returns all processes using ports in the given mode.
// On macOS: uses lsof. On Linux: tries lsof, falls back to /proc/net/tcp.
func Scan(mode ScanMode) ([]Process, error) {
	procs, err := scanLsof(mode)
	if err == nil {
		return procs, nil
	}
	// lsof failed — try proc fallback (Linux only)
	procs, procErr := scanProc(mode)
	if procErr == nil {
		return procs, nil
	}
	return nil, errors.New("portly requires lsof or /proc/net/tcp (Linux kernel 2.6+)")
}

// ScanPort returns processes using a specific port, or nil if free.
func ScanPort(port int) ([]Process, error) {
	all, err := Scan(AllConnections)
	if err != nil {
		return nil, err
	}
	var matches []Process
	for _, p := range all {
		if p.Port == port {
			matches = append(matches, p)
		}
	}
	return matches, nil
}

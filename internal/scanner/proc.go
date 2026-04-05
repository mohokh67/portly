package scanner

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type procEntry struct {
	Port    int
	Address string
	UID     int
	Inode   string
}

// scanProc reads /proc/net/tcp and /proc/net/tcp6 (Linux only).
func scanProc(mode ScanMode) ([]Process, error) {
	entries, err := readProcFile("/proc/net/tcp", mode)
	if err != nil {
		return nil, fmt.Errorf("cannot read /proc/net/tcp: %w", err)
	}
	entries6, _ := readProcFile("/proc/net/tcp6", mode) // best effort
	entries = append(entries, entries6...)

	inodeToPID := buildInodeMap()
	uidMap := buildUIDMap()

	var procs []Process
	for _, e := range entries {
		pid := inodeToPID[e.Inode]
		name := procName(pid)
		user, ok := uidMap[e.UID]
		if !ok {
			user = strconv.Itoa(e.UID)
		}
		procs = append(procs, Process{
			Port:    e.Port,
			Proto:   "TCP",
			PID:     pid,
			User:    user,
			Address: e.Address,
			Name:    name,
		})
	}
	return procs, nil
}

func readProcFile(path string, mode ScanMode) ([]procEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseProcNetTCP(string(data), mode)
}

// parseProcNetTCP parses the text content of /proc/net/tcp or /proc/net/tcp6.
// State 0A = LISTEN.
func parseProcNetTCP(content string, mode ScanMode) ([]procEntry, error) {
	var entries []procEntry
	lines := strings.Split(content, "\n")
	for _, line := range lines[1:] { // skip header
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		if mode == ListeningOnly && fields[3] != "0A" {
			continue
		}
		addr, port, err := parseHexAddrPort(fields[1])
		if err != nil {
			continue
		}
		uid, _ := strconv.Atoi(fields[7])
		inode := fields[9]
		entries = append(entries, procEntry{
			Port:    port,
			Address: addr,
			UID:     uid,
			Inode:   inode,
		})
	}
	return entries, nil
}

// parseHexAddrPort parses "AABBCCDD:PPPP" into IP string and port int.
// The IP is stored in little-endian hex on Linux.
func parseHexAddrPort(s string) (string, int, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid addr:port %q", s)
	}
	ipHex := parts[0]
	portHex := parts[1]

	portVal, err := strconv.ParseInt(portHex, 16, 32)
	if err != nil {
		return "", 0, err
	}

	ipBytes, err := hex.DecodeString(ipHex)
	if err != nil || len(ipBytes) < 4 {
		return "", 0, fmt.Errorf("bad ip hex %q", ipHex)
	}
	// Little-endian: reverse bytes for IPv4
	ip := fmt.Sprintf("%d.%d.%d.%d", ipBytes[3], ipBytes[2], ipBytes[1], ipBytes[0])
	return ip, int(portVal), nil
}

// buildInodeMap scans /proc/<pid>/fd/ symlinks to build inode → pid mapping.
func buildInodeMap() map[string]int {
	m := make(map[string]int)
	entries, err := filepath.Glob("/proc/[0-9]*/fd/*")
	if err != nil {
		return m
	}
	for _, fd := range entries {
		target, err := os.Readlink(fd)
		if err != nil {
			continue
		}
		if !strings.HasPrefix(target, "socket:[") {
			continue
		}
		inode := strings.TrimSuffix(strings.TrimPrefix(target, "socket:["), "]")
		parts := strings.Split(fd, "/")
		if len(parts) < 3 {
			continue
		}
		pid, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}
		m[inode] = pid
	}
	return m
}

// procName reads /proc/<pid>/comm for the process name.
func procName(pid int) string {
	if pid == 0 {
		return "unknown"
	}
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

// buildUIDMap reads /etc/passwd once and returns a uid → username map.
func buildUIDMap() map[int]string {
	m := make(map[int]string)
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return m
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Split(line, ":")
		if len(fields) < 3 {
			continue
		}
		uid, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		m[uid] = fields[0]
	}
	return m
}

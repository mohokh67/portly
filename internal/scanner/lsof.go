package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// scanLsof runs lsof and parses the output.
func scanLsof(mode ScanMode) ([]Process, error) {
	out, err := exec.Command("lsof", "-i", "-P", "-n").Output()
	if err != nil {
		// lsof not found or failed
		if _, ok := err.(*exec.Error); ok {
			return nil, fmt.Errorf("lsof not found")
		}
		// lsof exits 1 when it finds nothing — treat as empty, not an error
		if len(out) == 0 {
			return []Process{}, nil
		}
		// non-empty output with exit 1 still contains usable data — continue parsing
	}
	return parseLsof(string(out), mode)
}

// parseLsof parses the text output of `lsof -i -P -n`.
func parseLsof(output string, mode ScanMode) ([]Process, error) {
	var procs []Process
	lines := strings.Split(output, "\n")
	for _, line := range lines[1:] { // skip header
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		// Need at least: COMMAND PID USER FD TYPE DEVICE SIZE NODE NAME
		if len(fields) < 9 {
			continue
		}
		name := fields[0]
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		user := fields[2]

		// lsof output columns: COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME [STATE]
		// state "(LISTEN)" is a separate token after the address.
		state := ""
		nameField := strings.Join(fields[8:], " ")
		if last := fields[len(fields)-1]; strings.HasPrefix(last, "(") && strings.HasSuffix(last, ")") {
			state = strings.Trim(last, "()")
			nameField = strings.TrimSpace(strings.Join(fields[8:len(fields)-1], " "))
		}

		// Filter by mode
		isListen := state == "LISTEN"
		if mode == ListeningOnly && !isListen {
			continue
		}

		proto := "TCP"
		for _, f := range fields[:8] {
			if strings.EqualFold(f, "UDP") {
				proto = "UDP"
				break
			}
		}

		// Parse address and port from nameField
		// Could be "*:3000", "127.0.0.1:3000", "127.0.0.1:3000->127.0.0.1:54321"
		localPart := nameField
		if idx := strings.Index(nameField, "->"); idx != -1 {
			localPart = nameField[:idx]
		}

		lastColon := strings.LastIndex(localPart, ":")
		if lastColon == -1 {
			continue
		}
		addrStr := localPart[:lastColon]
		portStr := localPart[lastColon+1:]

		port, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		addr := addrStr
		if addr == "*" {
			addr = "0.0.0.0"
		}

		procs = append(procs, Process{
			Port:    port,
			Proto:   proto,
			PID:     pid,
			User:    user,
			Address: addr,
			Name:    name,
		})
	}
	return procs, nil
}

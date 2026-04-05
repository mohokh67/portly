package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mohokh67/portly/internal/scanner"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <port>",
	Short: "Show what is using a port (no kill prompt)",
	Args:  cobra.ExactArgs(1),
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	port, err := parsePort(args[0])
	if err != nil {
		return err
	}

	procs, err := scanner.ScanPort(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(procs) == 0 {
		fmt.Printf("port %d: free\n", port)
		return nil
	}

	for _, p := range procs {
		fmt.Printf("port %d: %s (PID %d, user: %s, proto: %s, addr: %s)\n",
			port, p.Name, p.PID, p.User, p.Proto, p.Address)
	}
	return nil
}

func parsePort(s string) (int, error) {
	port, err := strconv.Atoi(s)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("invalid port %q: must be an integer 1–65535", s)
	}
	return port, nil
}

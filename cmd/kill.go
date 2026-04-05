package cmd

import (
	"fmt"
	"os"

	"github.com/mohokh67/portly/internal/killer"
	"github.com/mohokh67/portly/internal/scanner"
	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill <port>",
	Short: "Kill all processes on a port immediately",
	Args:  cobra.ExactArgs(1),
	RunE:  runKill,
}

func init() {
	rootCmd.AddCommand(killCmd)
}

func runKill(cmd *cobra.Command, args []string) error {
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
		fmt.Printf("Nothing on port %d.\n", port)
		os.Exit(2)
	}

	for _, p := range procs {
		if err := killer.KillPID(p.PID); err != nil {
			fmt.Fprintf(os.Stderr, "failed to kill %s (PID %d): %v\n", p.Name, p.PID, err)
			os.Exit(1)
		}
		fmt.Printf("killed %s (PID %d) on :%d\n", p.Name, p.PID, port)
	}
	return nil
}

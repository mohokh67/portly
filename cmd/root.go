package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mohokh67/portly/internal/icons"
	"github.com/mohokh67/portly/internal/killer"
	"github.com/mohokh67/portly/internal/scanner"
	"github.com/mohokh67/portly/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "portly [port]",
	Short: "Manage ports interactively",
	Args:  cobra.ArbitraryArgs,
	RunE:  runRoot,
}

func Execute(version string) {
	rootCmd.Version = version
	rootCmd.PersistentFlags().String("icons", "auto", "icon style: nerdfont|emoji|none|auto")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		iconFlag, _ := cmd.Flags().GetString("icons")
		style := icons.ParseStyle(iconFlag)
		if style == icons.Auto {
			style = icons.DetectStyle()
		}
		return tui.Run(tui.Config{IconStyle: style})
	}

	if len(args) == 1 {
		port, err := parsePort(args[0])
		if err != nil {
			cmd.Usage()
			os.Exit(1)
		}
		return runPortLookup(port)
	}

	cmd.Usage()
	os.Exit(1)
	return nil
}

func runPortLookup(port int) error {
	procs, err := scanner.ScanPort(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(procs) == 0 {
		fmt.Printf("Port %d is free.\n", port)
		return nil
	}

	for _, p := range procs {
		fmt.Printf("port %d: %s (PID %d, user: %s)\n", port, p.Name, p.PID, p.User)
	}

	var prompt string
	if len(procs) == 1 {
		prompt = fmt.Sprintf("Kill %s (PID %d) on :%d?", procs[0].Name, procs[0].PID, port)
	} else {
		names := make([]string, len(procs))
		for i, p := range procs {
			names[i] = fmt.Sprintf("%s (PID %d)", p.Name, p.PID)
		}
		prompt = fmt.Sprintf("Kill %d processes on :%d (%s)?", len(procs), port, strings.Join(names, ", "))
	}
	if !promptYN(prompt) {
		return nil
	}

	for _, p := range procs {
		if err := killer.KillPID(p.PID); err != nil {
			fmt.Fprintf(os.Stderr, "failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("killed %s (PID %d)\n", p.Name, p.PID)
	}
	return nil
}

func promptYN(question string) bool {
	fmt.Printf("%s [y/N] ", question)
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	return ans == "y"
}

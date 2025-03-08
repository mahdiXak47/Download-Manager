package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "download-manager",
	Short: "A feature-rich download manager with TUI",
	Long: `A sophisticated download manager application written in Go with Terminal User Interface (TUI).
Features include multiple queues, download management, speed limiting, and scheduling capabilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This will be replaced with TUI initialization later
		fmt.Println("Download Manager - TUI coming soon!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
} 